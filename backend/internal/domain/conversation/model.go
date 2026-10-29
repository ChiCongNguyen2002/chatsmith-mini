package conversation

import "time"

type Conversation struct {
	ID        string    `gorm:"column:id;primaryKey;default:gen_random_uuid()"`
	UserID    string    `gorm:"column:user_id"`
	Title     string    `gorm:"column:title"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime:false;default:now()"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime:false;default:now()"`
}

func (Conversation) TableName() string {
	return tableConversations
}

type Page struct {
	Items      []Conversation
	NextCursor *string
}
