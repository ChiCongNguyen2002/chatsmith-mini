package chat

import (
	"encoding/json"
	"time"

	"chatsmith-mini/internal/domain/catalog"
	"chatsmith-mini/pkg/prompt"
)

type Message struct {
	ID             string          `gorm:"column:id;primaryKey;default:gen_random_uuid()"`
	ConversationID string          `gorm:"column:conversation_id"`
	Seq            int64           `gorm:"column:seq"`
	Role           Role            `gorm:"column:role"`
	Content        string          `gorm:"column:content"`
	Status         MessageStatus   `gorm:"column:status"`
	CreatedAt      time.Time       `gorm:"column:created_at;autoCreateTime:false;default:now()"`
	Model          *string         `gorm:"column:model"`
	InputTokens    *int64          `gorm:"column:input_tokens"`
	OutputTokens   *int64          `gorm:"column:output_tokens"`
	CostMicroUSD   *int64          `gorm:"column:cost_micro_usd"`
	UsageEstimated *bool           `gorm:"column:usage_estimated"`
	FirstTokenMS   *int            `gorm:"column:first_token_ms"`
	LatencyMS      *int            `gorm:"column:latency_ms"`
	Sources        json.RawMessage `gorm:"column:sources;->"`
}

func (Message) TableName() string {
	return tableMessages
}

type Usage struct {
	InputTokens  int64
	OutputTokens int64
	Estimated    bool
}

type Source struct {
	URL   string
	Title string
}

type Prior struct {
	Content   string
	Usage     Usage
	CostMicro int64
}

type Turn struct {
	ConversationID string
	UserMessage    Message
	Assistant      Message
	Model          catalog.Model
	Prompt         []prompt.Message
	RequestID      string
	Prior          *Prior
}

type FinishInput struct {
	ConversationID string
	MessageID      string
	Content        string
	Status         MessageStatus
	Usage          Usage
	CostMicroUSD   int64
	FirstTokenMS   *int
	LatencyMS      int
	Sources        []Source
}

type StartInput struct {
	UserID         string
	ConversationID string
	Content        string
	Model          string
	RequestID      string
}

type ContinueInput struct {
	UserID         string
	ConversationID string
	RequestID      string
}
