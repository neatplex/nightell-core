package models

import (
	"time"
)

type Story struct {
	ID            uint64    `gorm:"primaryKey" json:"-"`
	Identity      string    `gorm:"unique,size=60" json:"identity"`
	UserID        uint64    `json:"-"`
	User          User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	Caption       string    `gorm:"size:300" json:"caption"`
	AudioID       uint64    `json:"-"`
	Audio         File      `gorm:"foreignKey:AudioID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"audio"`
	ImageID       *uint64   `json:"-"`
	Image         *File     `gorm:"foreignKey:ImageID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"image"`
	IsPublished   bool      `gorm:"index" json:"is_published"`
	LikesCount    int       `json:"likes_count"`
	CommentsCount int       `json:"comments_count"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"-"`
}
