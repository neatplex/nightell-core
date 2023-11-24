package user

import (
	"errors"
	"fmt"
	"github.com/neatplex/nightel-core/internal/database"
	"github.com/neatplex/nightel-core/internal/models"
	"gorm.io/gorm"
)

type Service struct {
	database *database.Database
}

func (s *Service) FindByEmail(email string) (*models.User, error) {
	var user models.User
	r := s.database.Handler().Where("email = ?", email).First(&user)
	if r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.New(fmt.Sprintf("cannot query to find user by %s, err: %v", email, r.Error))
	}
	return &user, nil
}

func (s *Service) Create(user *models.User) error {
	r := s.database.Handler().Create(user)
	if r.Error != nil {
		return errors.New(fmt.Sprintf("cannot query to create user: %v", r.Error))
	}
	return nil
}

func New(database *database.Database) *Service {
	return &Service{database: database}
}
