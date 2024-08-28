package models

import (
	"time"
)

type Comment struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	UserId    uint64    `json:"user_id"`
	User      *User     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	PostId    uint64    `json:"post_id"`
	Post      *Post     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"post"`
	Text      string    `gorm:"size:255" json:"text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`
}
