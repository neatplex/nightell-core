package models

type File struct {
	ID        uint64  `gorm:"primaryKey" json:"-"`
	UserID    uint64  `json:"-"`
	User      User    `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Extension FileExt `gorm:"index" json:"extension"`
	Path      string  `gorm:"size=128" json:"path"`
}

type FileType string

const (
	FileTypeAudio FileType = "Audio"
	FileTypeImage FileType = "Image"
)

type FileExt string

func (fe FileExt) String() string {
	return string(fe)
}

const (
	FileExtMp3 FileExt = "MP3"
	FileExtJpg FileExt = "JPG"
)
