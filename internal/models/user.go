package models

import (
	"time"
)

type User struct {
	ID        uint64    `gorm:"primaryKey" json:"-"`
	Identity  string    `gorm:"unique" json:"identity"`
	Name      string    `json:"name"`
	Username  string    `gorm:"unique" json:"username"`
	Password  string    `json:"-"`
	Email     string    `gorm:"unique" json:"email"`
	IsTeller  bool      `gorm:"index" json:"is_teller"`
	IsBanned  bool      `gorm:"index" json:"is_banned"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`
}
