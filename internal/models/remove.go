package models

type Remove struct {
	Id     uint64 `gorm:"primaryKey" json:"id"`
	UserId uint64 `gorm:"unique" json:"user_id"`
	User   *User  `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	Code   string `gorm:"unique" json:"code"`
}
