package models

type Token struct {
	ID     uint64 `gorm:"primaryKey" json:"id"`
	UserID uint64 `json:"user_id"`
	User   *User  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	Value  string `gorm:"unique" json:"value"`
}
