package comment

import (
	"github.com/cockroachdb/errors"
	"github.com/neatplex/nightell-core/internal/database"
	"github.com/neatplex/nightell-core/internal/models"
	"gorm.io/gorm"
)

type Service struct {
	database *database.Database
}

func (s *Service) IndexByUser(userId uint64, lastId uint64, count int) ([]*models.Comment, error) {
	var comments []*models.Comment
	r := s.database.Handler().
		Where("user_id = ?", userId).
		Where("id < ? ORDER BY id DESC LIMIT ?", lastId, count).
		Preload("Post").
		Find(&comments)
	if r.Error != nil {
		return nil, errors.Wrapf(r.Error, "userId: %v", userId)
	}
	return comments, nil
}

func (s *Service) IndexByPost(postId uint64, lastId uint64, count int) ([]*models.Comment, error) {
	var comments []*models.Comment
	r := s.database.Handler().
		Where("post_id = ?", postId).
		Where("id < ? ORDER BY id DESC LIMIT ?", lastId, count).
		Preload("Post").
		Find(&comments)
	if r.Error != nil {
		return nil, errors.Wrapf(r.Error, "postId: %v", postId)
	}
	return comments, nil
}

func (s *Service) FindById(id uint64) (*models.Comment, error) {
	var comment models.Comment
	r := s.database.Handler().
		Where("id = ?", id).
		Preload("Post").
		Preload("User").
		Preload("User.Image").
		First(&comment)
	if r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.Wrapf(r.Error, "id: %v", id)
	}
	return &comment, nil
}

func (s *Service) Create(comment *models.Comment) (uint64, error) {
	r := s.database.Handler().Create(comment)
	if r.Error != nil {
		return 0, errors.Wrapf(r.Error, "comment: %v", comment)
	}
	return comment.ID, nil
}

func (s *Service) Delete(comment *models.Comment) {
	s.database.Handler().Delete(comment)
}

func New(database *database.Database) *Service {
	return &Service{database: database}
}
