package conversation

import "time"

type CreateRequest struct {
	Title string `json:"title"`
}

type RenameRequest struct {
	Title string `json:"title" binding:"required"`
}

type ListQuery struct {
	Limit  int    `form:"limit" binding:"min=0"`
	Cursor string `form:"cursor"`
}

type ConversationResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToResponse(item Conversation) ConversationResponse {
	return ConversationResponse{ID: item.ID, Title: item.Title, CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt}
}

func ToResponses(items []Conversation) []ConversationResponse {
	out := make([]ConversationResponse, 0, len(items))
	for _, item := range items {
		out = append(out, ToResponse(item))
	}

	return out
}
