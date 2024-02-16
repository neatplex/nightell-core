package like

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

func (s *Service) IndexByStoryIDWithUser(storyId uint64, lastId uint64, count int) ([]*models.Like, error) {
	var likes []*models.Like
	r := s.database.Handler().
		Preload("User").
		Where("story_id = ?", storyId).
		Where("id < ? ORDER BY id DESC LIMIT ?", lastId, count).
		Find(&likes)
	if r.Error != nil {
		return nil, fmt.Errorf("services: like: IndexByStoryIDWithUser: %s", r.Error)
	}
	return likes, nil
}

func (s *Service) Create(user *models.User, story *models.Story) (*models.Like, error) {
	var like models.Like
	r := s.database.Handler().FirstOrCreate(&like, &models.Like{UserID: user.ID, StoryID: story.ID})
	if r.Error != nil {
		return nil, fmt.Errorf("services: like: Create: %v, err: %v", like, r.Error)
	}
	return &like, r.Error
}

func (s *Service) FindById(id uint64) (*models.Like, error) {
	var model models.Like
	r := s.database.Handler().Where("id = ?", id).First(&model)
	if r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("services: like: FindById: `%d`, err: %v", id, r.Error)
	}
	return &model, nil
}

func (s *Service) Delete(id uint64) error {
	r := s.database.Handler().Delete(&models.Like{}, id)
	if r.Error != nil {
		return fmt.Errorf("services: like: Delete #%v, err: %v", id, r.Error)
	}
	return r.Error
}

func New(database *database.Database) *Service {
	return &Service{database: database}
}
