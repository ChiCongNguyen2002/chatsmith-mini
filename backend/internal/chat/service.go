package chat

import (
	"context"
	"strings"
	"time"
	"unicode/utf8"

	"chatsmith-mini/configs"
	"chatsmith-mini/constants"
	"chatsmith-mini/internal/catalog"
	"chatsmith-mini/internal/conversation"
	"chatsmith-mini/internal/provider"
	"chatsmith-mini/pkg/failure"
	"chatsmith-mini/pkg/logger"
)

const (
	finishTimeout     = 5 * time.Second
	logChatFinished   = "chat_finished"
	logChatFailed     = "chat_failed"
	logChatNotSaved   = "chat_result_not_saved"
	logKeyModel       = "model"
	logKeyStatus      = "status"
	logKeyInputTokens = "input_tokens"
	logKeyOutputToken = "output_tokens"
	logKeyCostMicro   = "cost_micro_usd"
	logKeyFirstToken  = "first_token_ms"
	logKeyLatency     = "latency_ms"
	logKeyEstimated   = "usage_estimated"
	logKeyConvID      = "conversation_id"
)

type Streamer interface {
	Stream(ctx context.Context, req provider.Request, onDelta func(string) error) (provider.Result, error)
}

type Service struct {
	repo          *Repo
	catalog       *catalog.Service
	conversations *conversation.Service
	streamers     map[constants.Provider]Streamer
	cfg           configs.ChatConfig
}

func newService(repo *Repo, cat *catalog.Service, conversations *conversation.Service, streamers map[constants.Provider]Streamer, cfg configs.ChatConfig) *Service {
	return &Service{repo: repo, catalog: cat, conversations: conversations, streamers: streamers, cfg: cfg}
}

func (s *Service) Start(ctx context.Context, in StartInput) (*Turn, error) {
	content := strings.TrimSpace(in.Content)
	if content == "" || utf8.RuneCountInString(content) > s.cfg.MaxMessageChars {
		return nil, failure.New(failure.CodeBadRequest)
	}

	model, err := s.catalog.Resolve(ctx, in.Model)
	if err != nil {
		return nil, err
	}
	if _, ok := s.streamers[model.Provider]; !ok {
		return nil, failure.New(failure.CodeModelUnavailable)
	}
	if _, err := s.conversations.Get(ctx, in.UserID, in.ConversationID); err != nil {
		return nil, err
	}
	if err := s.checkBudget(ctx); err != nil {
		return nil, err
	}

	staleBefore := time.Now().Add(-(s.cfg.Timeout + staleAfterMargin))
	user, assistant, err := s.repo.Begin(ctx, in.ConversationID, content, makeTitle(content), model.ID, staleBefore)
	if err != nil {
		return nil, err
	}

	history, err := s.repo.History(ctx, in.ConversationID, assistant.Seq)
	if err != nil {
		s.abandon(ctx, in.ConversationID, assistant.ID)
		return nil, err
	}

	return &Turn{
		ConversationID: in.ConversationID,
		UserMessage:    user,
		Assistant:      assistant,
		Model:          *model,
		Prompt:         buildPrompt(history, s.cfg.ContextTokens),
		RequestID:      in.RequestID,
	}, nil
}

func (s *Service) Run(ctx context.Context, turn *Turn, emit func(Event) error) {
	started := time.Now()
	emitErr := emit(Event{Name: EventCreated, Data: createdData{
		ConversationID:   turn.ConversationID,
		UserMessage:      messageRef{ID: turn.UserMessage.ID, Seq: turn.UserMessage.Seq},
		AssistantMessage: assistantRef{ID: turn.Assistant.ID, Seq: turn.Assistant.Seq, Model: turn.Model.ID},
	}})

	genCtx, cancel := context.WithTimeout(ctx, s.cfg.Timeout)
	defer cancel()

	var (
		content      strings.Builder
		firstTokenMS int
		result       provider.Result
		streamErr    error
	)
	if emitErr == nil {
		result, streamErr = s.streamers[turn.Model.Provider].Stream(genCtx, provider.Request{
			Model:            turn.Model.ID,
			Messages:         turn.Prompt,
			MaxOutputTokens:  s.cfg.MaxOutputTokens,
			SystemPrompt:     s.cfg.SystemPrompt,
			WebSearchEnabled: s.cfg.WebSearch,
			OnSearch: func(se provider.SearchEvent) error {
				switch se.Type {
				case "started":
					return emit(Event{Name: EventSearchStarted})
				case "source":
					return emit(Event{Name: EventSearchSource, Data: searchSourceData{URL: se.URL, Title: se.Title}})
				case "completed":
					return emit(Event{Name: EventSearchCompleted})
				}
				return nil
			},
		}, func(text string) error {
			if firstTokenMS == 0 {
				firstTokenMS = max(1, int(time.Since(started).Milliseconds()))
			}
			content.WriteString(text)
			emitErr = emit(Event{Name: EventDelta, Data: deltaData{Text: text}})

			return emitErr
		})
	}

	status, code := classify(ctx, genCtx, streamErr, emitErr)
	if status == conversation.StatusCompleted && result.Truncated {
		status = conversation.StatusTruncated
	}
	usage := usageFor(result, status, content.String(), turn.Prompt)
	costMicro := cost(usage, turn.Model.InputMicro, turn.Model.OutputMicro)
	latencyMS := int(time.Since(started).Milliseconds())

	var firstToken *int
	if firstTokenMS > 0 {
		firstToken = &firstTokenMS
	}

	saveCtx, saveCancel := context.WithTimeout(context.WithoutCancel(ctx), finishTimeout)
	defer saveCancel()
	saveErr := s.repo.Finish(saveCtx, finishInput{
		conversationID: turn.ConversationID,
		messageID:      turn.Assistant.ID,
		content:        content.String(),
		status:         status,
		usage:          usage,
		costMicroUSD:   costMicro,
		firstTokenMS:   firstToken,
		latencyMS:      latencyMS,
	})

	log := logger.FromContext(ctx)
	if streamErr != nil && status == conversation.StatusFailed {
		log.Error(logChatFailed, constants.LogKeyError, streamErr, logKeyModel, turn.Model.ID, logKeyConvID, turn.ConversationID)
	}
	if saveErr != nil {
		log.Error(logChatNotSaved, constants.LogKeyError, saveErr, logKeyConvID, turn.ConversationID)
		s.emitFailed(emit, turn, failure.CodeInternal)

		return
	}
	s.logFinished(log, turn, status, usage, costMicro, firstTokenMS, latencyMS)

	switch status {
	case conversation.StatusCompleted, conversation.StatusTruncated:
		_ = emit(Event{Name: EventCompleted, Data: completedData{
			MessageID:    turn.Assistant.ID,
			Status:       status,
			Model:        turn.Model.ID,
			Usage:        usage,
			CostMicroUSD: costMicro,
			FirstTokenMS: firstTokenMS,
			LatencyMS:    latencyMS,
		}})
	case conversation.StatusFailed:
		s.emitFailed(emit, turn, code)
	}
}

func (s *Service) checkBudget(ctx context.Context) error {
	budget := s.cfg.DailyBudgetMicro()
	if budget == 0 {
		return nil
	}

	now := time.Now().UTC()
	spent, err := s.repo.SpentSince(ctx, time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC))
	if err != nil {
		return err
	}
	if spent >= budget {
		return failure.New(failure.CodeBudgetExceeded)
	}

	return nil
}
