package chat

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"
	"unicode/utf8"

	"chatsmith-mini/configs"
	"chatsmith-mini/constants"
	"chatsmith-mini/internal/domain/catalog"
	"chatsmith-mini/internal/domain/conversation"
	"chatsmith-mini/pkg/failure"
	"chatsmith-mini/pkg/llm"
	"chatsmith-mini/pkg/logger"
	"chatsmith-mini/pkg/prompt"
)

type Streamer interface {
	Stream(ctx context.Context, req llm.Request, onDelta func(string) error) (llm.Result, error)
}

type Emit func(name string, payload any) error

type Service struct {
	repo          *Repo
	conversations *conversation.Service
	models        *catalog.Service
	streamers     map[constants.Provider]Streamer
	classifier    *llm.Classifier
	cfg           configs.ChatConfig
}

type reply struct {
	text         string
	firstTokenMS int
	result       llm.Result
	streamErr    error
	emitErr      error
	sources      []Source
}

type outcome struct {
	status    MessageStatus
	reason    failure.Code
	content   string
	usage     Usage
	costMicro int64
	latencyMS int
}

func newService(repo *Repo, conversations *conversation.Service, models *catalog.Service, streamers map[constants.Provider]Streamer, classifier *llm.Classifier, cfg configs.ChatConfig) *Service {
	return &Service{repo: repo, conversations: conversations, models: models, streamers: streamers, classifier: classifier, cfg: cfg}
}

func (s *Service) Start(ctx context.Context, in StartInput) (*Turn, error) {
	content := strings.TrimSpace(in.Content)
	if content == "" || utf8.RuneCountInString(content) > s.cfg.MaxMessageChars {
		return nil, failure.New(failure.CodeBadRequest)
	}

	model, err := s.models.Resolve(ctx, in.Model)
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
	user, assistant, err := s.repo.BeginTurn(ctx, in.ConversationID, content, prompt.Title(content), model.ID, staleBefore)
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
		Prompt:         prompt.Build(toPrompt(history), s.cfg.ContextTokens),
		RequestID:      in.RequestID,
	}, nil
}

func (s *Service) Continue(ctx context.Context, in ContinueInput) (*Turn, error) {
	if _, err := s.conversations.Get(ctx, in.UserID, in.ConversationID); err != nil {
		return nil, err
	}
	if err := s.checkBudget(ctx); err != nil {
		return nil, err
	}

	last, model, err := s.truncatedTail(ctx, in.ConversationID)
	if err != nil {
		return nil, err
	}

	staleBefore := time.Now().Add(-(s.cfg.Timeout + staleAfterMargin))
	assistant, err := s.repo.BeginContinue(ctx, last, staleBefore)
	if err != nil {
		return nil, err
	}

	history, err := s.repo.History(ctx, in.ConversationID, last.Seq)
	if err != nil {
		s.restoreTruncated(ctx, last)
		return nil, err
	}

	messages := append(prompt.Build(toPrompt(append(history, last)), s.cfg.ContextTokens), prompt.Message{Role: string(RoleUser), Content: ContinueInstruction})

	return &Turn{
		ConversationID: in.ConversationID,
		Assistant:      assistant,
		Model:          *model,
		Prompt:         messages,
		RequestID:      in.RequestID,
		Prior:          newPrior(last),
	}, nil
}

func (s *Service) ListMessages(ctx context.Context, userID, conversationID string, afterSeq int64, limit int) ([]Message, error) {
	if _, err := s.conversations.Get(ctx, userID, conversationID); err != nil {
		return nil, err
	}

	return s.repo.ListMessages(ctx, conversationID, afterSeq, clamp(limit, DefaultMessageLimit, MaxMessageLimit))
}

func (s *Service) Run(ctx context.Context, turn *Turn, emit Emit) {
	started := time.Now()
	genCtx, cancel := context.WithTimeout(ctx, s.cfg.Timeout)
	defer cancel()

	out := s.generate(genCtx, turn, emit, started)
	res := s.finalize(ctx, genCtx, turn, out, started)
	saveErr := s.save(ctx, turn, out, res)
	s.report(ctx, turn, emit, out, res, saveErr)
}

func (s *Service) generate(genCtx context.Context, turn *Turn, emit Emit, started time.Time) reply {
	if err := emit(EventCreated, newCreatedEvent(turn)); err != nil {
		return reply{emitErr: err}
	}

	var out reply
	var text strings.Builder
	webSearch := turn.Prior == nil && s.cfg.WebSearch && s.classifier.NeedsSearch(genCtx, turn.UserMessage.Content)
	out.result, out.streamErr = s.streamers[turn.Model.Provider].Stream(genCtx, llm.Request{
		Model:            turn.Model.ID,
		Messages:         toLLM(turn.Prompt),
		MaxOutputTokens:  s.cfg.MaxOutputTokens,
		SystemPrompt:     s.cfg.SystemPrompt,
		WebSearchEnabled: webSearch,
		OnSearch: func(se llm.SearchEvent) error {
			return s.onSearch(emit, &out, se)
		},
	}, func(delta string) error {
		if out.firstTokenMS == 0 {
			out.firstTokenMS = max(1, int(time.Since(started).Milliseconds()))
		}
		text.WriteString(delta)
		out.emitErr = emit(EventDelta, DeltaEvent{Text: delta})

		return out.emitErr
	})
	out.text = text.String()

	return out
}

func (s *Service) onSearch(emit Emit, out *reply, se llm.SearchEvent) error {
	switch se.Type {
	case "started":
		return emit(EventSearchStarted, nil)
	case "source":
		out.sources = append(out.sources, Source{URL: se.URL, Title: se.Title})

		return emit(EventSearchSource, SourceEvent{URL: se.URL, Title: se.Title})
	case "completed":
		return emit(EventSearchCompleted, nil)
	}

	return nil
}

func (s *Service) finalize(ctx, genCtx context.Context, turn *Turn, out reply, started time.Time) outcome {
	timedOut := errors.Is(genCtx.Err(), context.DeadlineExceeded)
	status, reason := classify(ctx.Err() != nil, timedOut, out.streamErr, out.emitErr)
	if status == StatusCompleted && out.result.Truncated {
		status = StatusTruncated
	}
	usage := usageFor(out.result, status, out.text, turn.Prompt)
	cost := costMicro(usage, turn.Model.InputMicro, turn.Model.OutputMicro)
	status, content, usage, cost := settle(turn.Prior, status, out.text, usage, cost)

	return outcome{
		status:    status,
		reason:    reason,
		content:   content,
		usage:     usage,
		costMicro: cost,
		latencyMS: int(time.Since(started).Milliseconds()),
	}
}

func (s *Service) save(ctx context.Context, turn *Turn, out reply, res outcome) error {
	var firstToken *int
	if out.firstTokenMS > 0 {
		firstToken = &out.firstTokenMS
	}
	saveCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), finishTimeout)
	defer cancel()

	return s.repo.FinishMessage(saveCtx, FinishInput{
		ConversationID: turn.ConversationID,
		MessageID:      turn.Assistant.ID,
		Content:        res.content,
		Status:         res.status,
		Usage:          res.usage,
		CostMicroUSD:   res.costMicro,
		FirstTokenMS:   firstToken,
		LatencyMS:      res.latencyMS,
		Sources:        out.sources,
	})
}

func (s *Service) report(ctx context.Context, turn *Turn, emit Emit, out reply, res outcome, saveErr error) {
	log := logger.FromContext(ctx)
	if out.streamErr != nil && res.status == StatusFailed {
		log.Error(logChatFailed, constants.LogKeyError, out.streamErr, logKeyModel, turn.Model.ID, logKeyConvID, turn.ConversationID)
	}
	if saveErr != nil {
		log.Error(logChatNotSaved, constants.LogKeyError, saveErr, logKeyConvID, turn.ConversationID)
		_ = emit(EventFailed, newFailedEvent(turn, failure.CodeInternal))

		return
	}
	logFinished(log, turn, res, out.firstTokenMS)

	switch res.status {
	case StatusCompleted, StatusTruncated:
		_ = emit(EventCompleted, newCompletedEvent(turn, res.status, res.usage, res.costMicro, out.firstTokenMS, res.latencyMS))
	case StatusFailed:
		_ = emit(EventFailed, newFailedEvent(turn, res.reason))
	}
}

func (s *Service) checkBudget(ctx context.Context) error {
	limit := s.cfg.DailyBudgetMicro()
	if limit == 0 {
		return nil
	}

	now := time.Now().UTC()
	spent, err := s.repo.SpentSince(ctx, time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC))
	if err != nil {
		return err
	}
	if spent >= limit {
		return failure.New(failure.CodeBudgetExceeded)
	}

	return nil
}

func (s *Service) truncatedTail(ctx context.Context, conversationID string) (Message, *catalog.Model, error) {
	last, err := s.repo.LatestMessage(ctx, conversationID)
	if err != nil {
		return last, nil, err
	}
	if last.Role != RoleAssistant || last.Status != StatusTruncated || last.Model == nil {
		return last, nil, failure.New(failure.CodeBadRequest)
	}

	model, err := s.models.Resolve(ctx, *last.Model)
	if err != nil {
		return last, nil, err
	}
	if _, ok := s.streamers[model.Provider]; !ok {
		return last, nil, failure.New(failure.CodeModelUnavailable)
	}

	return last, model, nil
}

func (s *Service) abandon(ctx context.Context, conversationID, messageID string) {
	s.abandonWith(ctx, FinishInput{ConversationID: conversationID, MessageID: messageID, Status: StatusFailed})
}

func (s *Service) restoreTruncated(ctx context.Context, last Message) {
	s.abandonWith(ctx, FinishInput{
		ConversationID: last.ConversationID,
		MessageID:      last.ID,
		Content:        last.Content,
		Status:         StatusTruncated,
		Usage:          Usage{InputTokens: derefInt64(last.InputTokens), OutputTokens: derefInt64(last.OutputTokens)},
		CostMicroUSD:   derefInt64(last.CostMicroUSD),
	})
}

func (s *Service) abandonWith(ctx context.Context, in FinishInput) {
	saveCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), finishTimeout)
	defer cancel()

	_ = s.repo.FinishMessage(saveCtx, in)
}

func classify(parentCanceled, timedOut bool, streamErr, emitErr error) (MessageStatus, failure.Code) {
	switch {
	case streamErr == nil && emitErr == nil:
		return StatusCompleted, ""
	case timedOut:
		return StatusFailed, failure.CodeProviderTimeout
	case emitErr != nil || parentCanceled:
		return StatusCanceled, ""
	default:
		return StatusFailed, failure.CodeProviderError
	}
}

func usageFor(result llm.Result, status MessageStatus, content string, messages []prompt.Message) Usage {
	usage := Usage{InputTokens: result.InputTokens, OutputTokens: result.OutputTokens}
	delivered := status == StatusCompleted || status == StatusTruncated
	if delivered || content == "" {
		return usage
	}

	usage.Estimated = true
	usage.OutputTokens = prompt.EstimateTokens(content)
	if usage.InputTokens == 0 {
		usage.InputTokens = prompt.Tokens(messages)
	}

	return usage
}

func settle(prior *Prior, status MessageStatus, content string, usage Usage, cost int64) (MessageStatus, string, Usage, int64) {
	if prior == nil {
		return status, content, usage, cost
	}
	if status == StatusFailed && content == "" {
		status = StatusTruncated
	}
	usage.InputTokens += prior.Usage.InputTokens
	usage.OutputTokens += prior.Usage.OutputTokens

	return status, prior.Content + content, usage, cost + prior.CostMicro
}

func costMicro(usage Usage, inputMicro, outputMicro int64) int64 {
	total := usage.InputTokens*inputMicro + usage.OutputTokens*outputMicro

	return (total + microPerMTok - 1) / microPerMTok
}

func newPrior(last Message) *Prior {
	return &Prior{
		Content:   last.Content,
		Usage:     Usage{InputTokens: derefInt64(last.InputTokens), OutputTokens: derefInt64(last.OutputTokens)},
		CostMicro: derefInt64(last.CostMicroUSD),
	}
}

func toPrompt(items []Message) []prompt.Message {
	out := make([]prompt.Message, 0, len(items))
	for _, item := range items {
		out = append(out, prompt.Message{Role: string(item.Role), Content: item.Content})
	}

	return out
}

func toLLM(items []prompt.Message) []llm.Message {
	out := make([]llm.Message, 0, len(items))
	for _, item := range items {
		out = append(out, llm.Message{Role: item.Role, Content: item.Content})
	}

	return out
}

func logFinished(log *slog.Logger, turn *Turn, res outcome, firstTokenMS int) {
	log.Info(logChatFinished,
		logKeyModel, turn.Model.ID,
		logKeyStatus, res.status,
		logKeyInputTokens, res.usage.InputTokens,
		logKeyOutputToken, res.usage.OutputTokens,
		logKeyEstimated, res.usage.Estimated,
		logKeyCostMicro, res.costMicro,
		logKeyFirstToken, firstTokenMS,
		logKeyLatency, res.latencyMS,
		logKeyConvID, turn.ConversationID,
	)
}

func derefInt64(v *int64) int64 {
	if v == nil {
		return 0
	}

	return *v
}

func clamp(limit, fallback, ceiling int) int {
	if limit <= 0 {
		return fallback
	}

	return min(limit, ceiling)
}
