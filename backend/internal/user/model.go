package user

import "time"

type User struct {
	ID        string    `gorm:"column:id;primaryKey;default:gen_random_uuid()" json:"id"`
	GoogleSub string    `gorm:"column:google_sub" json:"-"`
	Email     string    `gorm:"column:email" json:"email"`
	Name      string    `gorm:"column:name" json:"name"`
	Picture   string    `gorm:"column:picture" json:"picture"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime:false;default:now()" json:"created_at"`
}

func (User) TableName() string {
	return tableUsers
}

type Profile struct {
	GoogleSub string
	Email     string
	Name      string
	Picture   string
}
