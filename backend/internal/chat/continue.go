package chat

import (
	"context"
	"time"

	"chatsmith-mini/internal/conversation"
	"chatsmith-mini/internal/provider"
	"chatsmith-mini/pkg/failure"
)

type priorMessage struct {
	content   string
	usage     usageData
	costMicro int64
}

func newPrior(last conversation.Message) *priorMessage {
	return &priorMessage{
		content:   last.Content,
		usage:     usageData{InputTokens: derefInt64(last.InputTokens), OutputTokens: derefInt64(last.OutputTokens)},
		costMicro: derefInt64(last.CostMicroUSD),
	}
}

func (p *priorMessage) settle(status conversation.MessageStatus, content string, usage usageData, costMicro int64) (conversation.MessageStatus, string, usageData, int64) {
	if p == nil {
		return status, content, usage, costMicro
	}
	if status == conversation.StatusFailed && content == "" {
		status = conversation.StatusTruncated
	}
	usage.InputTokens += p.usage.InputTokens
	usage.OutputTokens += p.usage.OutputTokens

	return status, p.content + content, usage, costMicro + p.costMicro
}

func (s *Service) StartContinue(ctx context.Context, in ContinueInput) (*Turn, error) {
	if _, err := s.conversations.Get(ctx, in.UserID, in.ConversationID); err != nil {
		return nil, err
	}
	if err := s.checkBudget(ctx); err != nil {
		return nil, err
	}

	last, err := s.repo.Latest(ctx, in.ConversationID)
	if err != nil {
		return nil, err
	}
	if last.Role != conversation.RoleAssistant || last.Status != conversation.StatusTruncated || last.Model == nil {
		return nil, failure.New(failure.CodeBadRequest)
	}

	model, err := s.catalog.Resolve(ctx, *last.Model)
	if err != nil {
		return nil, err
	}
	if _, ok := s.streamers[model.Provider]; !ok {
		return nil, failure.New(failure.CodeModelUnavailable)
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

	prompt := append(buildPrompt(append(history, last), s.cfg.ContextTokens), provider.Message{Role: string(conversation.RoleUser), Content: continuePrompt})

	return &Turn{
		ConversationID: in.ConversationID,
		Assistant:      assistant,
		Model:          *model,
		Prompt:         prompt,
		RequestID:      in.RequestID,
		Prior:          newPrior(last),
	}, nil
}

func (s *Service) restoreTruncated(ctx context.Context, last conversation.Message) {
	s.abandonWith(ctx, finishInput{
		conversationID: last.ConversationID,
		messageID:      last.ID,
		content:        last.Content,
		status:         conversation.StatusTruncated,
		usage:          usageData{InputTokens: derefInt64(last.InputTokens), OutputTokens: derefInt64(last.OutputTokens)},
		costMicroUSD:   derefInt64(last.CostMicroUSD),
	})
}

func derefInt64(v *int64) int64 {
	if v == nil {
		return 0
	}

	return *v
}
