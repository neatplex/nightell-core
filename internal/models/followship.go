package models

import "time"

type Followship struct {
	ID         uint64    `gorm:"primaryKey" json:"id"`
	FollowerID uint64    `json:"follower_id"`
	Follower   *User     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"follower"`
	FolloweeID uint64    `json:"followee_id"`
	Followee   *User     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"followee"`
	CreatedAt  time.Time `json:"created_at"`
}
