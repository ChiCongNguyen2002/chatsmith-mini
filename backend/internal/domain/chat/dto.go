package chat

import (
	"encoding/json"
	"time"

	"chatsmith-mini/pkg/failure"
)

type SendRequest struct {
	Content string `json:"content"`
	Model   string `json:"model"`
}

type MessageListQuery struct {
	AfterSeq int64 `form:"after_seq" binding:"min=0"`
	Limit    int   `form:"limit" binding:"min=0"`
}

type MessageResponse struct {
	ID             string          `json:"id"`
	Seq            int64           `json:"seq"`
	Role           Role            `json:"role"`
	Content        string          `json:"content"`
	Status         MessageStatus   `json:"status"`
	CreatedAt      time.Time       `json:"created_at"`
	Model          *string         `json:"model"`
	InputTokens    *int64          `json:"input_tokens"`
	OutputTokens   *int64          `json:"output_tokens"`
	CostMicroUSD   *int64          `json:"cost_micro_usd"`
	UsageEstimated *bool           `json:"usage_estimated"`
	FirstTokenMS   *int            `json:"first_token_ms"`
	LatencyMS      *int            `json:"latency_ms"`
	Sources        json.RawMessage `json:"sources"`
}

type UsageResponse struct {
	InputTokens  int64 `json:"input_tokens"`
	OutputTokens int64 `json:"output_tokens"`
	Estimated    bool  `json:"estimated"`
}

type MessageRef struct {
	ID  string `json:"id"`
	Seq int64  `json:"seq"`
}

type AssistantRef struct {
	ID    string `json:"id"`
	Seq   int64  `json:"seq"`
	Model string `json:"model"`
}

type CreatedEvent struct {
	ConversationID   string       `json:"conversation_id"`
	UserMessage      MessageRef   `json:"user_message"`
	AssistantMessage AssistantRef `json:"assistant_message"`
}

type DeltaEvent struct {
	Text string `json:"text"`
}

type CompletedEvent struct {
	MessageID    string        `json:"message_id"`
	Status       MessageStatus `json:"status"`
	Model        string        `json:"model"`
	Usage        UsageResponse `json:"usage"`
	CostMicroUSD int64         `json:"cost_micro_usd"`
	FirstTokenMS int           `json:"first_token_ms"`
	LatencyMS    int           `json:"latency_ms"`
}

type FailedEvent struct {
	MessageID string        `json:"message_id"`
	Status    MessageStatus `json:"status"`
	Code      failure.Code  `json:"code"`
	RequestID string        `json:"request_id"`
}

type SourceEvent struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

func ToMessageResponses(items []Message) []MessageResponse {
	out := make([]MessageResponse, 0, len(items))
	for _, item := range items {
		out = append(out, MessageResponse{
			ID:             item.ID,
			Seq:            item.Seq,
			Role:           item.Role,
			Content:        item.Content,
			Status:         item.Status,
			CreatedAt:      item.CreatedAt,
			Model:          item.Model,
			InputTokens:    item.InputTokens,
			OutputTokens:   item.OutputTokens,
			CostMicroUSD:   item.CostMicroUSD,
			UsageEstimated: item.UsageEstimated,
			FirstTokenMS:   item.FirstTokenMS,
			LatencyMS:      item.LatencyMS,
			Sources:        item.Sources,
		})
	}

	return out
}

func newUsageResponse(u Usage) UsageResponse {
	return UsageResponse{InputTokens: u.InputTokens, OutputTokens: u.OutputTokens, Estimated: u.Estimated}
}

func newCreatedEvent(turn *Turn) CreatedEvent {
	return CreatedEvent{
		ConversationID:   turn.ConversationID,
		UserMessage:      MessageRef{ID: turn.UserMessage.ID, Seq: turn.UserMessage.Seq},
		AssistantMessage: AssistantRef{ID: turn.Assistant.ID, Seq: turn.Assistant.Seq, Model: turn.Model.ID},
	}
}

func newCompletedEvent(turn *Turn, status MessageStatus, usage Usage, costMicro int64, firstTokenMS, latencyMS int) CompletedEvent {
	return CompletedEvent{
		MessageID:    turn.Assistant.ID,
		Status:       status,
		Model:        turn.Model.ID,
		Usage:        newUsageResponse(usage),
		CostMicroUSD: costMicro,
		FirstTokenMS: firstTokenMS,
		LatencyMS:    latencyMS,
	}
}

func newFailedEvent(turn *Turn, code failure.Code) FailedEvent {
	return FailedEvent{MessageID: turn.Assistant.ID, Status: StatusFailed, Code: code, RequestID: turn.RequestID}
}

type SourceJSON struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}
