package token

import (
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/neatplex/nightel-core/internal/database"
	"github.com/neatplex/nightel-core/internal/models"
	"gorm.io/gorm"
)

type Service struct {
	database *database.Database
}

func (s *Service) Create(user *models.User) (string, error) {
	value := uuid.NewString()
	r := s.database.Handler().Create(&models.Token{
		UserID: user.ID,
		Value:  value,
	})
	return value, errors.Wrapf(r.Error, "user: %v", user.ID)
}

func (s *Service) FindOrCreate(user *models.User) (string, error) {
	var token models.Token
	r := s.database.Handler().Where("user_id = ?", user.ID).First(&token)
	if r.Error != nil && errors.Is(r.Error, gorm.ErrRecordNotFound) {
		tokenValue, err := s.Create(user)
		return tokenValue, errors.WithStack(err)
	}
	return token.Value, errors.Wrapf(r.Error, "user: %v", user.ID)
}

func (s *Service) FindByValue(value string) (*models.Token, error) {
	var token models.Token
	r := s.database.Handler().Where("value = ?", value).Preload("User").First(&token)
	if r.Error != nil && errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &token, errors.Wrapf(r.Error, "value: %v", value)
}

func New(database *database.Database) *Service {
	return &Service{database: database}
}
