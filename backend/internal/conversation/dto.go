package conversation

type ConversationPage struct {
	Items      []Conversation `json:"items"`
	NextCursor *string        `json:"next_cursor"`
}
