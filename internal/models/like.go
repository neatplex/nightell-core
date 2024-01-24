package models

import (
	"time"
)

type Like struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	UserID    uint64    `json:"user_id"`
	User      *User     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	StoryID   uint64    `json:"story_id"`
	Story     *Story    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"story"`
	CreatedAt time.Time `json:"created_at"`
}
