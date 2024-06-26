package user

import (
	"github.com/cockroachdb/errors"
	"github.com/neatplex/nightell-core/internal/database"
	"github.com/neatplex/nightell-core/internal/mailer"
	"github.com/neatplex/nightell-core/internal/models"
	"gorm.io/gorm"
)

type Service struct {
	database *database.Database
	mailer   *mailer.Mailer
}

func (s *Service) FindBy(field string, value interface{}) (*models.User, error) {
	var user models.User
	r := s.database.Handler().Preload("Image").Where(field+" = ?", value).First(&user)
	if r.Error != nil && errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, errors.Wrapf(r.Error, "field: %v, value: %v", field, value)
}

func (s *Service) UpdateName(user *models.User, name string) (*models.User, error) {
	user.Name = name
	r := s.database.Handler().Save(user)
	return user, errors.Wrapf(r.Error, "user: %v", user)
}

func (s *Service) UpdateImage(user *models.User, imageID uint64) (*models.User, error) {
	user.ImageID = &imageID
	r := s.database.Handler().Save(user)
	return user, errors.Wrapf(r.Error, "user: %v", user)
}

func (s *Service) UpdateBio(user *models.User, bio string) (*models.User, error) {
	user.Bio = bio
	r := s.database.Handler().Save(user)
	return user, errors.Wrapf(r.Error, "user: %v", user)
}

func (s *Service) UpdateUsername(user *models.User, username string) (*models.User, error) {
	oldUser, _ := s.FindBy("username", username)
	if oldUser != nil && oldUser.ID != user.ID {
		return nil, ErrUsernameAlreadyExist
	}

	user.Username = username
	r := s.database.Handler().Save(user)
	return user, errors.Wrapf(r.Error, "user: %v", user)
}

func (s *Service) Delete(user *models.User) error {
	r := s.database.Handler().Delete(user)
	return errors.Wrapf(r.Error, "user: %v", user)
}

func (s *Service) Create(user *models.User) error {
	defer func() {
		go s.mailer.SendWellcome(user.Email, user.Username)
	}()

	r := s.database.Handler().Create(user)
	return errors.Wrapf(r.Error, "user: %v", user)
}

func New(database *database.Database, mailer *mailer.Mailer) *Service {
	return &Service{database: database, mailer: mailer}
}
