package models

import (
	"time"
)

type Story struct {
	ID            uint64    `gorm:"primaryKey" json:"id"`
	UserID        uint64    `json:"user_id"`
	User          *User     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	Caption       string    `gorm:"size:300" json:"caption"`
	AudioID       uint64    `json:"audio_id"`
	Audio         *File     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"audio"`
	ImageID       *uint64   `json:"image_id"`
	Image         *File     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"image"`
	Likes         []*Like   `json:"likes"`
	LikesCount    uint64    `json:"likes_count"`
	CommentsCount uint64    `json:"comments_count"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"-"`
}
