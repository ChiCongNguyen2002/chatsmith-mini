package chat

import "time"

const (
	EventCreated         = "message.created"
	EventDelta           = "message.delta"
	EventCompleted       = "message.completed"
	EventFailed          = "message.failed"
	EventSearchStarted   = "search.started"
	EventSearchSource    = "search.source"
	EventSearchCompleted = "search.completed"
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

const (
	microPerMTok     = 1_000_000
	charsPerToken    = 4
	tokenOverhead    = 4
	maxHistoryRows   = 200
	titleRunes       = 40
	staleAfterMargin = 30 * time.Second
	uniqueViolation  = "23505"
)

const (
	whereStaleStreaming   = "conversation_id = ? AND status = ? AND created_at < ?"
	whereConversationRows = "conversation_id = ?"
	whereHistory          = "conversation_id = ? AND status = ? AND seq < ?"
	whereSpent            = "role = ? AND created_at >= ?"
	selectMaxSeq          = "COALESCE(MAX(seq), 0)"
	selectSumCost         = "COALESCE(SUM(cost_micro_usd), 0)"
	orderNewestSeqFirst   = "seq DESC"
	columnStatus          = "status"
	columnRole            = "role"
	columnContent         = "content"
)

const (
	touchConversationSQL = `UPDATE conversations
SET updated_at = now(), title = CASE WHEN title = '' THEN ? ELSE title END
WHERE id = ?`

	finishMessageSQL = `UPDATE messages
SET content = ?, status = ?, input_tokens = ?, output_tokens = ?, cost_micro_usd = ?,
    usage_estimated = ?, first_token_ms = ?, latency_ms = ?
WHERE id = ?`

	touchOnlySQL = `UPDATE conversations SET updated_at = now() WHERE id = ?`
)
