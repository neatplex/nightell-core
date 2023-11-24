package models

type Token struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint
	User   User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Value  string `gorm:"unique"`
}
