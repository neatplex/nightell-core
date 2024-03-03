package like

import (
	"github.com/cockroachdb/errors"
	"github.com/neatplex/nightel-core/internal/database"
	"github.com/neatplex/nightel-core/internal/models"
	"gorm.io/gorm"
)

type Service struct {
	database *database.Database
}

func (s *Service) IndexByPostIDWithUser(postId uint64, lastId uint64, count int) ([]*models.Like, error) {
	var likes []*models.Like
	if count > 100 {
		count = 100
	}
	r := s.database.Handler().
		Preload("User").
		Where("post_id = ?", postId).
		Where("id < ? ORDER BY id DESC LIMIT ?", lastId, count).
		Find(&likes)
	return likes, errors.Wrapf(r.Error, "postId: %v, lastId: %v, count: %v", postId, lastId, count)
}

func (s *Service) Create(user *models.User, post *models.Post) (*models.Like, error) {
	var like models.Like
	r := s.database.Handler().FirstOrCreate(&like, &models.Like{UserID: user.ID, PostID: post.ID})
	return &like, errors.Wrapf(r.Error, "user: %v, post: %v", user, post)
}

func (s *Service) FindById(id uint64) (*models.Like, error) {
	var model models.Like
	r := s.database.Handler().Where("id = ?", id).First(&model)
	if r.Error != nil && errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &model, errors.Wrapf(r.Error, "id: %v", id)
}

func (s *Service) Delete(id uint64) error {
	r := s.database.Handler().Delete(&models.Like{}, id)
	return errors.Wrapf(r.Error, "id: %v", id)
}

func New(database *database.Database) *Service {
	return &Service{database: database}
}
