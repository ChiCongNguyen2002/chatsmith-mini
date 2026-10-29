package user

import "time"

type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Picture   string    `json:"picture"`
	CreatedAt time.Time `json:"created_at"`
}

func ToResponse(item User) UserResponse {
	return UserResponse{ID: item.ID, Email: item.Email, Name: item.Name, Picture: item.Picture, CreatedAt: item.CreatedAt}
}
