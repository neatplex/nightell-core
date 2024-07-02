package remove

import (
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/neatplex/nightell-core/internal/database"
	"github.com/neatplex/nightell-core/internal/models"
	"gorm.io/gorm"
)

type Service struct {
	database *database.Database
}

func (s *Service) Create(user *models.User) (string, error) {
	value := uuid.NewString()
	r := s.database.Handler().Create(&models.Remove{
		UserID: user.ID,
		Code:   value,
	})
	return value, errors.Wrapf(r.Error, "user: %v", user.ID)
}

func (s *Service) FindOrCreate(user *models.User) (string, error) {
	var model models.Remove
	r := s.database.Handler().Where("user_id = ?", user.ID).First(&model)
	if r.Error != nil && errors.Is(r.Error, gorm.ErrRecordNotFound) {
		code, err := s.Create(user)
		return code, errors.WithStack(err)
	}
	return model.Code, errors.Wrapf(r.Error, "user: %v", user.ID)
}

func (s *Service) FindBy(field string, value interface{}) (*models.Remove, error) {
	var model models.Remove
	r := s.database.Handler().Preload("User").Where(field+" = ?", value).First(&model)
	if r.Error != nil && errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &model, errors.Wrapf(r.Error, "field: %v, value: %v", field, value)
}

func New(database *database.Database) *Service {
	return &Service{database: database}
}
