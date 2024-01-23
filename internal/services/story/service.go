package story

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/neatplex/nightel-core/internal/database"
	"github.com/neatplex/nightel-core/internal/models"
	"gorm.io/gorm"
	"strings"
)

type Service struct {
	database *database.Database
}

func (s *Service) Index(userId uint64) ([]models.Story, error) {
	var stories []models.Story
	r := s.database.Handler().
		Where("user_id = ?", userId).
		Preload("Audio").Preload("Image").
		Find(&stories)
	if r.Error != nil {
		return nil, fmt.Errorf("services: story: Index: %v", r.Error)
	}
	return stories, nil
}

func (s *Service) Feed(lastId uint64) ([]models.Story, error) {
	var stories []models.Story
	r := s.database.Handler().
		Where("SELECT * FROM app.stories WHERE id > ? ORDER BY id LIMIT 2", lastId).
		Preload("Audio").
		Preload("Image").
		Find(&stories)
	if r.Error != nil {
		return nil, fmt.Errorf("services: story: Feed: %s", r.Error)
	}
	return stories, nil
}

func (s *Service) Create(story *models.Story) (string, error) {
	if story.Identity == "" {
		story.Identity = strings.ReplaceAll(uuid.NewString(), "-", "0")
	}

	r := s.database.Handler().Create(story)
	if r.Error != nil {
		return "", fmt.Errorf("services: story: Create: %v", r.Error)
	}
	return story.Identity, nil
}

func (s *Service) UpdateCaption(story *models.Story, caption string) *models.Story {
	story.Caption = caption
	s.database.Handler().Save(story)
	return story
}

func (s *Service) Delete(story *models.Story) {
	s.database.Handler().Delete(story)
}

func (s *Service) FindByIdentity(identity string) (*models.Story, error) {
	var story models.Story
	r := s.database.Handler().
		Where("identity = ?", identity).
		Preload("Audio").Preload("Image").
		First(&story)
	if r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("services: story: FindByIdentity: `%s`, err: %v", identity, r.Error)
	}
	return &story, nil
}

func New(database *database.Database) *Service {
	return &Service{database: database}
}
