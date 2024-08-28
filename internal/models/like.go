package models

import (
	"time"
)

type Like struct {
	Id        uint64    `gorm:"primaryKey" json:"id"`
	UserId    uint64    `json:"user_id"`
	User      *User     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	PostId    uint64    `json:"post_id"`
	Post      *Post     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"post"`
	CreatedAt time.Time `json:"created_at"`
}
