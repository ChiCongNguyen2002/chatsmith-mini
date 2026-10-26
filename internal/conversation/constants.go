package conversation

const (
	MaxTitleRunes       = 200
	DefaultListLimit    = 20
	MaxListLimit        = 100
	DefaultMessageLimit = 50
	MaxMessageLimit     = 200
)

const (
	tableUsers         = "users"
	tableConversations = "conversations"
	tableMessages      = "messages"
)

const (
	whereOwnedConversation = "id = ? AND user_id = ?"
	whereUserConversations = "user_id = ?"
	whereBeforeCursor      = "(updated_at, id) < (?, ?)"
	orderNewestFirst       = "updated_at DESC, id DESC"
	whereMessagesAfterSeq  = "conversation_id = ? AND seq > ?"
	orderOldestFirst       = "seq ASC"
)

const cursorSeparator = "_"
