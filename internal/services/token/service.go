package token

import (
	"errors"
	"fmt"
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
	if r.Error != nil {
		return "", errors.New(fmt.Sprintf("cannot query to create token for %d, err: %v", user.ID, r.Error))
	}
	return value, r.Error
}

func (s *Service) FindOrCreate(user *models.User) (string, error) {
	var token models.Token
	r := s.database.Handler().Where("user_id = ?", user.ID).First(&token)
	if r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			return s.Create(user)
		}
		return "", errors.New(fmt.Sprintf("cannot query to find token by user %d, err: %v", user.ID, r.Error))
	}
	return token.Value, nil
}

func (s *Service) FindByValue(value string) (*models.Token, error) {
	var token models.Token
	r := s.database.Handler().Where("value = ?", value).Preload("User").First(&token)
	if r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.New(fmt.Sprintf("cannot query to find token by %s, err: %v", value, r.Error))
	}
	return &token, nil
}

func New(database *database.Database) *Service {
	return &Service{database: database}
}
