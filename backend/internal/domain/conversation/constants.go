package conversation

const (
	tableConversations = "conversations"

	MaxTitleRunes    = 200
	DefaultListLimit = 20
	MaxListLimit     = 100

	whereOwned             = "id = ? AND user_id = ?"
	whereUserConversations = "user_id = ?"
	whereBeforeCursor      = "(updated_at, id) < (?, ?)"
	orderNewestFirst       = "updated_at DESC, id DESC"
)
