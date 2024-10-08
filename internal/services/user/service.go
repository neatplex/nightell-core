package user

import (
	"github.com/cockroachdb/errors"
	"github.com/neatplex/nightell-core/internal/database"
	"github.com/neatplex/nightell-core/internal/mailer"
	"github.com/neatplex/nightell-core/internal/models"
	"golang.org/x/crypto/bcrypt"
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

func (s *Service) CheckPassword(user *models.User, password string) bool {
	return user != nil && bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil
}

func (s *Service) Search(q string, lastId uint64, count int) ([]*models.User, error) {
	var users []*models.User
	r := s.database.Handler().
		Where("(name LIKE ? OR username LIKE ?)", "%"+q+"%", "%"+q+"%").
		Where("id < ? ORDER BY id DESC LIMIT ?", lastId, count).
		Preload("Image").
		Find(&users)
	if r.Error != nil {
		return nil, errors.Wrapf(r.Error, "q: %v, lastId: %v, count: %v", q, lastId, count)
	}
	return users, nil
}

func (s *Service) UpdateName(user *models.User, name string) (*models.User, error) {
	user.Name = name
	r := s.database.Handler().Save(user)
	return user, errors.Wrapf(r.Error, "user: %v", user)
}

func (s *Service) UpdateEmail(user *models.User, email string) (*models.User, error) {
	oldUser, _ := s.FindBy("email", email)
	if oldUser != nil && oldUser.Id != user.Id {
		return nil, ErrEmailAlreadyExist
	}

	user.Email = email
	r := s.database.Handler().Save(user)
	return user, errors.Wrapf(r.Error, "user: %v", user)
}

func (s *Service) UpdatePassword(user *models.User, password string) (*models.User, error) {
	user.Password = password
	r := s.database.Handler().Save(user)
	return user, errors.Wrapf(r.Error, "user: %v", user)
}

func (s *Service) UpdateImage(user *models.User, imageID uint64) (*models.User, error) {
	user.ImageId = &imageID
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
	if oldUser != nil && oldUser.Id != user.Id {
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
		go s.mailer.SendWelcome(user.Email, user.Username)
	}()

	r := s.database.Handler().Create(user)
	return errors.Wrapf(r.Error, "user: %v", user)
}

func (s *Service) FindByEmailOrCreate(email string) (*models.User, error) {
	user, err := s.FindBy("email", email)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if user != nil {
		return user, nil
	}

	user = &models.User{
		Username: email,
		Email:    email,
		IsBanned: false,
		Password: "",
	}

	return user, errors.WithStack(s.Create(user))
}

func New(database *database.Database, mailer *mailer.Mailer) *Service {
	return &Service{database: database, mailer: mailer}
}
