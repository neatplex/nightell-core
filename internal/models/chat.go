package models

import "time"

type Chat struct {
	ID        uint64     `gorm:"primaryKey" json:"id"`
	FromID    uint64     `json:"from_id"`
	From      *User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"from"`
	ToID      uint64     `json:"to_id"`
	To        *User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"to"`
	Messages  []*Message `json:"messages"`
	UpdatedAt time.Time  `gorm:"index" json:"updated_at"`
	CreatedAt time.Time  `gorm:"index" json:"created_at"`
}
