package models

type Token struct {
	ID     uint64 `gorm:"primaryKey" json:"-"`
	UserID uint64 `json:"user_id"`
	User   User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	Value  string `gorm:"unique" json:"value"`
}
