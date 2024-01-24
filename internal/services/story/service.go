package story

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

func (s *Service) Index(userId uint64) ([]*models.Story, error) {
	var stories []*models.Story
	r := s.database.Handler().
		Where("user_id = ?", userId).
		Preload("Audio").Preload("Image").
		Find(&stories)
	if r.Error != nil {
		return nil, fmt.Errorf("services: story: Index: %v", r.Error)
	}
	return stories, nil
}

func (s *Service) Feed(userId uint64, lastId uint64) ([]*models.Story, error) {
	var stories []*models.Story
	r := s.database.Handler().
		Where("id < ? ORDER BY id DESC LIMIT 3", lastId).
		Preload("Audio").
		Preload("Image").
		Find(&stories)
	if r.Error != nil {
		return nil, fmt.Errorf("services: story: Feed: %s", r.Error)
	}
	for _, story := range stories {
		if err := s.attachLikes(userId, story); err != nil {
			return nil, err
		}
	}
	return stories, nil
}

func (s *Service) attachLikes(userId uint64, story *models.Story) error {
	var likes []*models.Like

	r := s.database.Handler().
		Where("story_id = ?", story.ID).
		Where("user_id = ?", userId).
		Preload("User").
		Find(&likes)
	if r.Error != nil {
		return fmt.Errorf("services: story: attachLikes: %s", r.Error)
	}
	story.Likes = likes

	var count int64
	r = s.database.Handler().Model(&models.Like{}).Where("story_id = ?", story.ID).Count(&count)
	if r.Error != nil {
		return fmt.Errorf("services: story: attachLikes: %s", r.Error)
	}
	story.LikesCount = uint64(count)

	return nil
}

func (s *Service) Create(story *models.Story) (uint64, error) {
	r := s.database.Handler().Create(story)
	if r.Error != nil {
		return 0, fmt.Errorf("services: story: Create: %v", r.Error)
	}
	return story.ID, nil
}

func (s *Service) UpdateCaption(story *models.Story, caption string) *models.Story {
	story.Caption = caption
	s.database.Handler().Save(story)
	return story
}

func (s *Service) Delete(story *models.Story) {
	s.database.Handler().Delete(story)
}

func (s *Service) FindById(id uint64) (*models.Story, error) {
	var story models.Story
	r := s.database.Handler().
		Where("id = ?", id).
		Preload("Audio").Preload("Image").
		First(&story)
	if r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("services: story: FindById: `%d`, err: %v", id, r.Error)
	}
	return &story, nil
}

func New(database *database.Database) *Service {
	return &Service{database: database}
}
