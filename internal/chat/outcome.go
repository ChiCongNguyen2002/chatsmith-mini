package chat

import (
	"context"
	"errors"
	"log/slog"

	"chatsmith-mini/internal/conversation"
	"chatsmith-mini/internal/provider"
	"chatsmith-mini/pkg/failure"
)

func classify(ctx, genCtx context.Context, streamErr, emitErr error) (conversation.MessageStatus, failure.Code) {
	switch {
	case streamErr == nil && emitErr == nil:
		return conversation.StatusCompleted, ""
	case errors.Is(genCtx.Err(), context.DeadlineExceeded):
		return conversation.StatusFailed, failure.CodeProviderTimeout
	case emitErr != nil || ctx.Err() != nil:
		return conversation.StatusCanceled, ""
	default:
		return conversation.StatusFailed, failure.CodeProviderError
	}
}

func usageFor(result provider.Result, status conversation.MessageStatus, content string, prompt []provider.Message) usageData {
	usage := usageData{InputTokens: result.InputTokens, OutputTokens: result.OutputTokens}
	if status == conversation.StatusCompleted || content == "" {
		return usage
	}

	usage.Estimated = true
	usage.OutputTokens = estimateTokens(content)
	if usage.InputTokens == 0 {
		usage.InputTokens = promptTokens(prompt)
	}

	return usage
}

func (s *Service) abandon(ctx context.Context, conversationID, messageID string) {
	saveCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), finishTimeout)
	defer cancel()

	_ = s.repo.Finish(saveCtx, finishInput{conversationID: conversationID, messageID: messageID, status: conversation.StatusFailed})
}

func (s *Service) emitFailed(emit func(Event) error, turn *Turn, code failure.Code) {
	_ = emit(Event{Name: EventFailed, Data: failedData{
		MessageID: turn.Assistant.ID,
		Status:    conversation.StatusFailed,
		Code:      code,
		RequestID: turn.RequestID,
	}})
}

func (s *Service) logFinished(log *slog.Logger, turn *Turn, status conversation.MessageStatus, usage usageData, costMicro int64, firstTokenMS, latencyMS int) {
	log.Info(logChatFinished,
		logKeyModel, turn.Model.ID,
		logKeyStatus, status,
		logKeyInputTokens, usage.InputTokens,
		logKeyOutputToken, usage.OutputTokens,
		logKeyEstimated, usage.Estimated,
		logKeyCostMicro, costMicro,
		logKeyFirstToken, firstTokenMS,
		logKeyLatency, latencyMS,
		logKeyConvID, turn.ConversationID,
	)
}
