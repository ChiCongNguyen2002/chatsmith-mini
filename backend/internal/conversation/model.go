package conversation

import "time"

type Conversation struct {
	ID        string    `gorm:"column:id;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    string    `gorm:"column:user_id" json:"-"`
	Title     string    `gorm:"column:title" json:"title"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime:false;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime:false;default:now()" json:"updated_at"`
}

func (Conversation) TableName() string {
	return tableConversations
}

type Message struct {
	ID             string        `gorm:"column:id;primaryKey;default:gen_random_uuid()" json:"id"`
	ConversationID string        `gorm:"column:conversation_id" json:"-"`
	Seq            int64         `gorm:"column:seq" json:"seq"`
	Role           Role          `gorm:"column:role" json:"role"`
	Content        string        `gorm:"column:content" json:"content"`
	Status         MessageStatus `gorm:"column:status" json:"status"`
	CreatedAt      time.Time     `gorm:"column:created_at;autoCreateTime:false;default:now()" json:"created_at"`
	Model          *string       `gorm:"column:model" json:"model"`
	InputTokens    *int64        `gorm:"column:input_tokens" json:"input_tokens"`
	OutputTokens   *int64        `gorm:"column:output_tokens" json:"output_tokens"`
	CostMicroUSD   *int64        `gorm:"column:cost_micro_usd" json:"cost_micro_usd"`
	UsageEstimated *bool         `gorm:"column:usage_estimated" json:"usage_estimated"`
	FirstTokenMS   *int          `gorm:"column:first_token_ms" json:"first_token_ms"`
	LatencyMS      *int          `gorm:"column:latency_ms" json:"latency_ms"`
}

func (Message) TableName() string {
	return tableMessages
}
