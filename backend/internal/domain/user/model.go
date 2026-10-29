package user

import (
	"time"

	"chatsmith-mini/constants"
)

type User struct {
	ID        string    `gorm:"column:id;primaryKey;default:gen_random_uuid()"`
	Email     string    `gorm:"column:email"`
	Name      string    `gorm:"column:name"`
	Picture   string    `gorm:"column:picture"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime:false;default:now()"`
}

func (User) TableName() string {
	return tableUsers
}

type Identity struct {
	Provider      constants.AuthProvider
	Subject       string
	Email         string
	EmailVerified bool
	Name          string
	Picture       string
}

type identityRow struct {
	Provider      string    `gorm:"column:provider;primaryKey"`
	Subject       string    `gorm:"column:subject;primaryKey"`
	UserID        string    `gorm:"column:user_id"`
	Email         string    `gorm:"column:email"`
	EmailVerified bool      `gorm:"column:email_verified"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime:false;default:now()"`
}

func (identityRow) TableName() string {
	return tableIdentities
}
