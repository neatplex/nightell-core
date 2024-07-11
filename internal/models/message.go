package models

import (
	"time"
)

type Message struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	ChatID    uint64    `json:"chat_id"`
	Chat      *User     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"chat"`
	UserID    uint64    `json:"user_id"`
	User      *User     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user"`
	Content   string    `json:"content"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}
