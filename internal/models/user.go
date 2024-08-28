package models

import (
	"time"
)

type User struct {
	Id        uint64    `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:255" json:"name"`
	Username  string    `gorm:"unique" json:"username"`
	Password  string    `gorm:"size:255" json:"-"`
	Email     string    `gorm:"unique" json:"email"`
	Bio       string    `gorm:"size:255" json:"bio"`
	IsBanned  bool      `gorm:"index" json:"-"`
	ImageId   *uint64   `json:"image_id"`
	Image     *File     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"image"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`
}
