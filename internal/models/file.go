package models

import "time"

type File struct {
	Id        uint64    `gorm:"primaryKey" json:"id"`
	UserId    uint64    `json:"user_id"`
	User      *User     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	Extension Extension `gorm:"index" json:"extension"`
	Path      string    `gorm:"size:255" json:"path"`
	CreatedAt time.Time `json:"created_at"`
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
