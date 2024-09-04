package models

type Token struct {
	Id     uint64 `gorm:"primaryKey" json:"id"`
	UserId uint64 `json:"user_id"`
	User   *User  `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	Value  string `gorm:"unique" json:"value"`
}
