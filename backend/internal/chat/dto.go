package chat

import (
	"chatsmith-mini/internal/catalog"
	"chatsmith-mini/internal/conversation"
	"chatsmith-mini/internal/provider"
	"chatsmith-mini/pkg/failure"
)

type StartInput struct {
	UserID         string
	ConversationID string
	Content        string
	Model          string
	RequestID      string
}

type Turn struct {
	ConversationID string
	UserMessage    conversation.Message
	Assistant      conversation.Message
	Model          catalog.Model
	Prompt         []provider.Message
	RequestID      string
}

type Event struct {
	Name string
	Data any
}

type messageRef struct {
	ID  string `json:"id"`
	Seq int64  `json:"seq"`
}

type assistantRef struct {
	ID    string `json:"id"`
	Seq   int64  `json:"seq"`
	Model string `json:"model"`
}

type createdData struct {
	ConversationID   string       `json:"conversation_id"`
	UserMessage      messageRef   `json:"user_message"`
	AssistantMessage assistantRef `json:"assistant_message"`
}

type deltaData struct {
	Text string `json:"text"`
}

type usageData struct {
	InputTokens  int64 `json:"input_tokens"`
	OutputTokens int64 `json:"output_tokens"`
	Estimated    bool  `json:"estimated"`
}

type completedData struct {
	MessageID    string                     `json:"message_id"`
	Status       conversation.MessageStatus `json:"status"`
	Model        string                     `json:"model"`
	Usage        usageData                  `json:"usage"`
	CostMicroUSD int64                      `json:"cost_micro_usd"`
	FirstTokenMS int                        `json:"first_token_ms"`
	LatencyMS    int                        `json:"latency_ms"`
}

type failedData struct {
	MessageID string                     `json:"message_id"`
	Status    conversation.MessageStatus `json:"status"`
	Code      failure.Code               `json:"code"`
	RequestID string                     `json:"request_id"`
}
