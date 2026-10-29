package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type Session struct {
	TokenHash string    `gorm:"column:token_hash;primaryKey"`
	UserID    string    `gorm:"column:user_id"`
	ExpiresAt time.Time `gorm:"column:expires_at"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime:false;default:now()"`
}

func (Session) TableName() string {
	return tableSessions
}

func (s Session) Expired(now time.Time) bool {
	return !now.Before(s.ExpiresAt)
}

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))

	return hex.EncodeToString(sum[:])
}
