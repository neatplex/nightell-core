package file

import (
	"errors"
	"fmt"
	"github.com/neatplex/nightel-core/internal/database"
	"github.com/neatplex/nightel-core/internal/models"
	"github.com/neatplex/nightel-core/internal/s3"
	"gorm.io/gorm"
	"io"
)

type Service struct {
	database *database.Database
	s3       *s3.S3
}

func (s *Service) ToFileType(typeName string) (models.FileType, error) {
	switch typeName {
	case "mp3":
		return models.FileTypeMp3, nil
	}
	return "", errors.New("invalid file type")
}

func (s *Service) Download(path string) ([]byte, error) {
	return s.s3.Get(path)
}

func (s *Service) Upload(reader io.Reader) (string, error) {
	path := "random.mp3"
	return path, s.s3.Put(path, reader)
}

func (s *Service) Create(file *models.File) error {
	r := s.database.Handler().Create(file)
	if r.Error != nil {
		return errors.New(fmt.Sprintf("cannot query to create file: %v", r.Error))
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
		return nil, errors.New(fmt.Sprintf("cannot query to find file by %s, err: %v", path, r.Error))
	}
	return &file, nil
}

func New(database *database.Database, s3 *s3.S3) *Service {
	return &Service{database: database, s3: s3}
}
