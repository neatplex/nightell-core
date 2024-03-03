package models

import (
	"time"
)

type User struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Username  string    `gorm:"unique" json:"username"`
	Password  string    `json:"-"`
	Email     string    `gorm:"unique" json:"email"`
	Bio       string    `json:"bio"`
	IsBanned  bool      `gorm:"index" json:"-"`
	ImageID   *uint64   `json:"image_id"`
	Image     *File     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"image"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`
}
