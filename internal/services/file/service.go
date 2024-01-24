package file

import (
	"errors"
	"github.com/google/uuid"
	"github.com/neatplex/nightel-core/internal/database"
	"github.com/neatplex/nightel-core/internal/models"
	"github.com/neatplex/nightel-core/internal/s3"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
	"io"
	"time"
)

type Service struct {
	database *database.Database
	s3       *s3.S3
}

func (s *Service) ToExtension(name string) (models.Extension, error) {
	switch name {
	case "MP3":
		return models.FileExtMp3, nil
	case "JPG":
		return models.FileExtJpg, nil
	}
	return "", eris.Errorf("cannot convert `%s` to models.Extension", name)
}

func (s *Service) ToType(extension models.Extension) (models.FileType, error) {
	switch extension {
	case models.FileExtMp3:
		return models.FileTypeAudio, nil
	case models.FileExtJpg:
		return models.FileTypeImage, nil
	}
	return "", eris.Errorf("cannot convert `%v` to models.FileType", extension)
}

func (s *Service) Download(path string) ([]byte, error) {
	return s.s3.Get(path)
}

func (s *Service) Upload(reader io.Reader, extension models.Extension) (string, error) {
	path := time.Now().Format("2006/01/02/") + uuid.NewString() + "." + extension.String()
	return path, s.s3.Put(path, reader)
}

func (s *Service) Create(file *models.File) error {
	r := s.database.Handler().Create(file)
	if r.Error != nil {
		return eris.Errorf("cannot query to create file: %v", r.Error)
	}
	return nil
}

func (s *Service) FindByPath(path string) (*models.File, error) {
	var file models.File
	r := s.database.Handler().Where("path = ?", path).First(&file)
	if r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, eris.Errorf("cannot query to find file by %s, err: %v", path, r.Error)
	}
	return &file, nil
}

func New(database *database.Database, s3 *s3.S3) *Service {
	return &Service{database: database, s3: s3}
}
