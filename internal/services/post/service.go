package post

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

func (s *Service) Index(userId uint64) ([]*models.Post, error) {
	var posts []*models.Post
	r := s.database.Handler().
		Where("user_id = ?", userId).
		Preload("Audio").Preload("Image").
		Find(&posts)
	if r.Error != nil {
		return nil, fmt.Errorf("services: post: Index: %v", r.Error)
	}
	for _, post := range posts {
		if err := s.attachLikes(userId, post); err != nil {
			return nil, err
		}
	}
	return posts, nil
}

func (s *Service) Feed(userId uint64, lastId uint64, count int) ([]*models.Post, error) {
	var posts []*models.Post
	r := s.database.Handler().
		Where("id < ? ORDER BY id DESC LIMIT ?", lastId, count).
		Preload("Audio").
		Preload("Image").
		Find(&posts)
	if r.Error != nil {
		return nil, fmt.Errorf("services: post: Feed: %s", r.Error)
	}
	for _, post := range posts {
		if err := s.attachLikes(userId, post); err != nil {
			return nil, err
		}
	}
	return posts, nil
}

func (s *Service) attachLikes(userId uint64, post *models.Post) error {
	var likes []*models.Like

	r := s.database.Handler().
		Where("post_id = ?", post.ID).
		Where("user_id = ?", userId).
		Preload("User").
		Find(&likes)
	if r.Error != nil {
		return fmt.Errorf("services: post: attachLikes: %s", r.Error)
	}
	post.Likes = likes

	var count int64
	r = s.database.Handler().Model(&models.Like{}).Where("post_id = ?", post.ID).Count(&count)
	if r.Error != nil {
		return fmt.Errorf("services: post: attachLikes: %s", r.Error)
	}
	post.LikesCount = uint64(count)

	return nil
}

func (s *Service) Create(post *models.Post) (uint64, error) {
	r := s.database.Handler().Create(post)
	if r.Error != nil {
		return 0, fmt.Errorf("services: post: Create: %v", r.Error)
	}
	return post.ID, nil
}

func (s *Service) UpdateCaption(post *models.Post, caption string) *models.Post {
	post.Caption = caption
	s.database.Handler().Save(post)
	return post
}

func (s *Service) Delete(post *models.Post) {
	s.database.Handler().Delete(post)
}

func (s *Service) FindById(id uint64) (*models.Post, error) {
	var post models.Post
	r := s.database.Handler().
		Where("id = ?", id).
		Preload("Audio").Preload("Image").
		First(&post)
	if r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("services: post: FindById: `%d`, err: %v", id, r.Error)
	}
	if err := s.attachLikes(post.UserID, &post); err != nil {
		return nil, err
	}
	return &post, nil
}

func (s *Service) FindBy(field string, value interface{}) (*models.Post, error) {
	var post models.Post
	r := s.database.Handler().Where(field+" = ?", value).First(&post)
	if r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("services: post: FindBy%s: %v, err: %v", field, value, r.Error)
	}
	return &post, nil
}

func New(database *database.Database) *Service {
	return &Service{database: database}
}
