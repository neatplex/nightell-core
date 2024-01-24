package models

type File struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	UserID    uint64    `json:"user_id"`
	User      *User     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	Extension Extension `gorm:"index" json:"extension"`
	Path      string    `gorm:"size=128" json:"path"`
}

type FileType string

const (
	FileTypeAudio FileType = "Audio"
	FileTypeImage FileType = "Image"
)

type Extension string

func (fe Extension) String() string {
	return string(fe)
}

const (
	FileExtMp3 Extension = "MP3"
	FileExtJpg Extension = "JPG"
)
