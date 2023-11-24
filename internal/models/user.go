package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Username  string    `gorm:"unique" json:"username"`
	Password  string    `json:"-"`
	Email     string    `gorm:"unique" json:"email"`
	Status    Status    `gorm:"index" json:"status"`
	IsTeller  bool      `gorm:"index" json:"is_teller"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`
}

type Status string

const (
	StatusRegistered Status = "Registered"
	StatusDeleted    Status = "Deleted"
)
