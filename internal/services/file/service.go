package file

import (
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/neatplex/nightell-core/internal/database"
	"github.com/neatplex/nightell-core/internal/models"
	"github.com/neatplex/nightell-core/internal/s3"
	"gorm.io/gorm"
	"io"
	"strings"
	"time"
)

type Service struct {
	database *database.Database
	s3       *s3.S3
}

func (s *Service) ExtensionFromString(name string) (models.Extension, error) {
	switch name {
	case "MP3":
		return models.FileExtMp3, nil
	case "JPG":
		return models.FileExtJpg, nil
	}
	return "", errors.Errorf("cannot create extension from `%s`", name)
}

func (s *Service) TypeFromExtension(extension models.Extension) (models.FileType, error) {
	switch extension {
	case models.FileExtMp3:
		return models.FileTypeAudio, nil
	case models.FileExtJpg:
		return models.FileTypeImage, nil
	}
	return "", errors.Errorf("cannot create type from extension `%s`", extension.String())
}

func (s *Service) Download(path string) ([]byte, error) {
	file, err := s.s3.Get(path)
	return file, errors.Wrapf(err, "path: `%v`", path)
}

func (s *Service) Upload(reader io.Reader, path string) error {
	return errors.WithStack(s.s3.Put(path, reader))
}

func (s *Service) Path(extension models.Extension) string {
	return time.Now().Format("2006/01/02/") + uuid.NewString() + "." + strings.ToLower(extension.String())
}

func (s *Service) Create(file *models.File) error {
	r := s.database.Handler().Create(file)
	return errors.Wrapf(r.Error, "file: `%v`", file)
}

func (s *Service) FindByID(id uint64) (*models.File, error) {
	var file models.File
	r := s.database.Handler().First(&file, id)
	if r.Error != nil && errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &file, errors.Wrapf(r.Error, "id: `%v`", id)
}

func (s *Service) FindByPath(path string) (*models.File, error) {
	var file models.File
	r := s.database.Handler().Where("path = ?", path).First(&file)
	if r.Error != nil && errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &file, errors.Wrapf(r.Error, "path: `%v`", path)
}

func New(database *database.Database, s3 *s3.S3) *Service {
	return &Service{database: database, s3: s3}
}
