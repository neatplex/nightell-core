package message

import (
	"github.com/cockroachdb/errors"
	"github.com/neatplex/nightell-core/internal/database"
	"github.com/neatplex/nightell-core/internal/models"
	"gorm.io/gorm"
)

type Service struct {
	database *database.Database
}

func (s *Service) OneByChat(chatId uint64) (*models.Message, error) {
	var message *models.Message
	r := s.database.Handler().
		Where("chat_id = ?", chatId).
		Order("id DESC").
		Limit(1).
		First(message)
	if r.Error != nil && errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return message, errors.Wrapf(r.Error, "chatId: %v", chatId)
}

func (s *Service) Create(model *models.Message) (uint64, error) {
	r := s.database.Handler().Create(model)
	return model.ID, errors.Wrapf(r.Error, "model: %v", model)
}

func (s *Service) Delete(model *models.Message) {
	s.database.Handler().Delete(model)
}

func New(database *database.Database) *Service {
	return &Service{database: database}
}
