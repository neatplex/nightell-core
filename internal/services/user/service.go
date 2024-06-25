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

func (s *Service) UpdateName(user *models.User, name string) *models.User {
	user.Name = name
	s.database.Handler().Save(user)
	return user
}

func (s *Service) UpdateImage(user *models.User, imageID uint64) *models.User {
	user.ImageID = &imageID
	s.database.Handler().Save(user)
	return user
}

func (s *Service) UpdateBio(user *models.User, bio string) *models.User {
	user.Bio = bio
	s.database.Handler().Save(user)
	return user
}

func (s *Service) UpdateUsername(user *models.User, username string) (*models.User, error) {
	oldUser, _ := s.FindBy("username", username)
	if oldUser != nil && oldUser.ID != user.ID {
		return nil, ErrUsernameAlreadyExist
	}

	user.Username = username
	s.database.Handler().Save(user)
	return user, nil
}

func (s *Service) Create(user *models.User) error {
	defer func() {
		s.mailer.SendWellcome(user.Email, user.Username)
	}()

	r := s.database.Handler().Create(user)
	return errors.Wrapf(r.Error, "user: %v", user)
}

func New(database *database.Database, mailer *mailer.Mailer) *Service {
	return &Service{database: database, mailer: mailer}
}
