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
)

const (
	tableMessages = "messages"

	whereMessagesAfterSeq = "conversation_id = ? AND seq > ?"
	orderOldestFirst      = "seq ASC"
	whereStaleStreaming   = "conversation_id = ? AND status = ? AND created_at < ?"
	whereConversationRows = "conversation_id = ?"
	whereHistory          = "conversation_id = ? AND status IN ? AND seq < ?"
	whereSpent            = "role = ? AND created_at >= ?"
	selectMaxSeq          = "COALESCE(MAX(seq), 0)"
	selectSumCost         = "COALESCE(SUM(cost_micro_usd), 0)"
	orderNewestSeqFirst   = "seq DESC"
	columnStatus          = "status"
	columnRole            = "role"
	columnContent         = "content"

	touchConversationSQL = `UPDATE conversations
SET updated_at = now(), title = CASE WHEN title = '' THEN ? ELSE title END
WHERE id = ?`

	finishMessageSQL = `UPDATE messages
SET content = ?, status = ?, input_tokens = ?, output_tokens = ?, cost_micro_usd = ?,
    usage_estimated = ?, first_token_ms = ?, latency_ms = ?, sources = COALESCE(?, sources)
WHERE id = ?`

	touchOnlySQL    = `UPDATE conversations SET updated_at = now() WHERE id = ?`
	updateStatusSQL = `UPDATE messages SET status = ? WHERE id = ? AND status = ?`
)
