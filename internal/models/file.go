package models

type File struct {
	ID     uint64   `gorm:"primaryKey" json:"-"`
	UserID uint64   `json:"-"`
	User   User     `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Type   FileType `gorm:"index,size=5" json:"type"`
	Path   string   `gorm:"size=128" json:"path"`
}

type FileType string

const (
	FileTypeMp3 FileType = "mp3"
	FileTypeJpg FileType = "jpg"
)
