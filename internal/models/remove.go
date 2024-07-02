package models

type Remove struct {
	ID     uint64 `gorm:"primaryKey" json:"id"`
	UserID uint64 `gorm:"unique" json:"user_id"`
	User   *User  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	Code   string `gorm:"unique" json:"code"`
}
