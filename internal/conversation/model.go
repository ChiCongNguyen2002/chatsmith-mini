package conversation

import "time"

type User struct {
	ID        string    `gorm:"column:id;primaryKey"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime:false;default:now()"`
}

func (User) TableName() string {
	return tableUsers
}

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
	ID             string    `gorm:"column:id;primaryKey" json:"id"`
	ConversationID string    `gorm:"column:conversation_id" json:"-"`
	Seq            int64     `gorm:"column:seq" json:"seq"`
	Role           string    `gorm:"column:role" json:"role"`
	Content        string    `gorm:"column:content" json:"content"`
	Status         string    `gorm:"column:status" json:"status"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime:false" json:"created_at"`
}

func (Message) TableName() string {
	return tableMessages
}
