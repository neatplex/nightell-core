package models

import (
	"time"
)

type Story struct {
	ID            uint      `gorm:"primaryKey" json:"-"`
	Identity      string    `gorm:"unique" json:"identity"`
	UserID        uint      `json:"-"`
	User          User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	Caption       string    `gorm:"type:text" json:"caption"`
	Audio         string    `json:"audio"`
	Image         string    `json:"image"`
	IsPublished   bool      `gorm:"index" json:"is_published"`
	LikesCount    int       `json:"likes_count"`
	CommentsCount int       `json:"comments_count"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"-"`
}
