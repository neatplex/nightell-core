package models

import "time"

type Followship struct {
	Id         uint64    `gorm:"primaryKey" json:"id"`
	FollowerId uint64    `gorm:"index:idx_follower_id_followee_id,priority:1" json:"follower_id"`
	Follower   *User     `gorm:"foreignKey:FollowerId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"follower"`
	FolloweeId uint64    `gorm:"index:idx_follower_id_followee_id,priority:2" json:"followee_id"`
	Followee   *User     `gorm:"foreignKey:FolloweeId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"followee"`
	CreatedAt  time.Time `json:"created_at"`
}
