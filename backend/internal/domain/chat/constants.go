package chat

import "time"

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

type MessageStatus string

const (
	StatusStreaming MessageStatus = "streaming"
	StatusCompleted MessageStatus = "completed"
	StatusFailed    MessageStatus = "failed"
	StatusCanceled  MessageStatus = "canceled"
	StatusTruncated MessageStatus = "truncated"
)

const (
	EventCreated         = "message.created"
	EventDelta           = "message.delta"
	EventCompleted       = "message.completed"
	EventFailed          = "message.failed"
	EventSearchStarted   = "search.started"
	EventSearchSource    = "search.source"
	EventSearchCompleted = "search.completed"

	ContinueInstruction = "Tiếp tục đúng từ chỗ bạn vừa dừng. Không lặp lại phần đã viết và không mở đầu lại."
	microPerMTok        = 1_000_000
	DefaultMessageLimit = 50
	MaxMessageLimit     = 200
	finishTimeout       = 5 * time.Second
	staleAfterMargin    = 30 * time.Second
	maxHistoryRows      = 200
	uniqueViolation     = "23505"

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

	chatMeterName      = "chatsmith/chat"
	metricChatOutcome  = "chat.outcome"
	metricChatDuration = "chat.duration"
	attrStatus         = "status"
	attrReason         = "reason"
	attrModel          = "model"
	statusRejected     = "rejected"
	modelUnknown       = "unknown"
)

var chatDurationBuckets = []float64{0.25, 0.5, 1, 1.5, 2, 3, 5, 10, 20, 60, 120}
