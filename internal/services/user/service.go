package user

import (
	"github.com/cockroachdb/errors"
	"github.com/neatplex/nightell-core/internal/database"
	"github.com/neatplex/nightell-core/internal/models"
	"gorm.io/gorm"
)

type Service struct {
	database *database.Database
}

func (s *Service) FindById(id uint64) (*models.User, error) {
	var user models.User
	r := s.database.Handler().Where("id = ?", id).First(&user)
	if r.Error != nil && errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, errors.Wrapf(r.Error, "id: %v", id)
}

func (s *Service) FindByEmail(email string) (*models.User, error) {
	var user models.User
	r := s.database.Handler().Where("email = ?", email).First(&user)
	if r.Error != nil && errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, errors.Wrapf(r.Error, "email: %v", email)
}

func (s *Service) FindByUsername(email string) (*models.User, error) {
	var user models.User
	r := s.database.Handler().Where("username = ?", email).First(&user)
	if r.Error != nil && errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, errors.Wrapf(r.Error, "username: %v", email)
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
	return errors.Wrapf(r.Error, "user: %v", user)
}

func New(database *database.Database) *Service {
	return &Service{database: database}
}
