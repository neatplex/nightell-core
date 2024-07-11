package chat

import (
	"github.com/cockroachdb/errors"
	"github.com/neatplex/nightell-core/internal/database"
	"github.com/neatplex/nightell-core/internal/models"
	"github.com/neatplex/nightell-core/internal/services/message"
	"gorm.io/gorm"
)

type Service struct {
	database       *database.Database
	messageService *message.Service
}

func (s *Service) ByUser(userId uint64, lastId uint64, count int) ([]*models.Chat, error) {
	var chats []*models.Chat
	if count > 100 {
		count = 100
	}
	r := s.database.Handler().
		Where("(to_id = ? OR from_id = ?)", userId, userId).
		Where("id < ? ORDER BY updated_at DESC LIMIT ?", lastId, count).
		Preload("From").
		Preload("From.Image").
		Preload("To").
		Preload("To.Image").
		Find(&chats)
	if r.Error != nil {
		return nil, errors.Wrapf(r.Error, "userId: %v, lastId: %v, count: %v", userId, lastId, count)
	}
	for _, chat := range chats {
		m, err := s.messageService.OneByChat(chat.ID)
		if err != nil {
			return nil, errors.Wrapf(r.Error, "userId: %v, lastId: %v, count: %v", userId, lastId, count)
		}
		if m != nil {
			chat.Messages = []*models.Message{m}
		}
	}
	return chats, nil
}

func (s *Service) OneByUsers(users []uint64) (*models.Chat, error) {
	var chat *models.Chat
	r := s.database.Handler().
		Where("(to_id in (?) OR from_id in (?))", users, users).
		Preload("From").
		Preload("From.Image").
		Preload("To").
		Preload("To.Image").
		First(&chat)
	if r.Error != nil && errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return chat, errors.Wrapf(r.Error, "users: %v", users)
}

func (s *Service) OneBy(field string, value interface{}) (*models.Chat, error) {
	var model models.Chat
	r := s.database.Handler().
		Where(field+" = ?", value).
		Preload("Audio").
		Preload("Image").
		Preload("User").
		Preload("User.Image").
		First(&model)
	if r.Error != nil && errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &model, errors.Wrapf(r.Error, "field: %v, value: %v", field, value)
}

func (s *Service) Create(model *models.Chat) error {
	r := s.database.Handler().Create(model)
	return errors.Wrapf(r.Error, "model: %v", model)
}

func (s *Service) Delete(model *models.Chat) {
	s.database.Handler().Delete(model)
}

func New(database *database.Database, ms *message.Service) *Service {
	return &Service{database: database, messageService: ms}
}
