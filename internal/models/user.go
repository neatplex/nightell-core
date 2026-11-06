package models

import (
	"time"
)

type User struct {
	Id        uint64    `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:255" json:"name"`
	Username  string    `gorm:"unique" json:"username"`
	Password  string    `gorm:"size:255" json:"-"`
	Email     string    `gorm:"unique" json:"-"`
	Bio       string    `gorm:"size:255" json:"bio"`
	IsBanned  bool      `gorm:"index" json:"-"`
	ImageId   *uint64   `json:"image_id"`
	Image     *File     `gorm:"foreignKey:ImageId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"image"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`
}

// PrivateUser is a User as its owner sees it, including the email.
type PrivateUser struct {
	*User
	Email string `json:"email"`
}

// Private returns the user with the fields only its owner may see.
func (u *User) Private() *PrivateUser {
	if u == nil {
		return nil
	}
	return &PrivateUser{User: u, Email: u.Email}
}
