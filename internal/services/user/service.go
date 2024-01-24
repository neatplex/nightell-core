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
		return nil, fmt.Errorf("cannot query to find user by email: %s, err: %v", email, r.Error)
	}
	return &user, nil
}

func (s *Service) FindByUsername(email string) (*models.User, error) {
	var user models.User
	r := s.database.Handler().Where("username = ?", email).First(&user)
	if r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("cannot query to find user by username: %s, err: %v", email, r.Error)
	}
	return &user, nil
}

func (s *Service) UpdateName(user *models.User, name string) *models.User {
	user.Name = name
	s.database.Handler().Save(user)
	return user
}

func (s *Service) UpdateBio(user *models.User, bio string) *models.User {
	user.Bio = bio
	s.database.Handler().Save(user)
	return user
}

func (s *Service) UpdateUsername(user *models.User, username string) (*models.User, error) {
	oldUser, _ := s.FindByUsername(username)
	if oldUser != nil && oldUser.ID != user.ID {
		return nil, ErrUsernameAlreadyExist
	}

	user.Username = username
	s.database.Handler().Save(user)
	return user, nil
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
