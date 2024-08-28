package post

import (
	"github.com/cockroachdb/errors"
	"github.com/neatplex/nightell-core/internal/database"
	"github.com/neatplex/nightell-core/internal/models"
	"gorm.io/gorm"
)

type Service struct {
	database *database.Database
}

func (s *Service) Index(userId uint64, lastId uint64, count int) ([]*models.Post, error) {
	var posts []*models.Post
	r := s.database.Handler().
		Where("user_id = ?", userId).
		Where("id < ? ORDER BY id DESC LIMIT ?", lastId, count).
		Preload("Audio").
		Preload("Image").
		Preload("User").
		Preload("User.Image").
		Find(&posts)
	if r.Error != nil {
		return nil, errors.Wrapf(r.Error, "userId: %v", userId)
	}
	for _, post := range posts {
		if err := s.attachRelations(userId, post); err != nil {
			return nil, errors.Wrapf(err, "userId: %v, postId: %v", userId, post.Id)
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
		Preload("User").
		Preload("User.Image").
		Find(&posts)
	if r.Error != nil {
		return nil, errors.Wrapf(r.Error, "userId: %v, lastId: %v, count: %v", userId, lastId, count)
	}
	for _, post := range posts {
		if err := s.attachRelations(userId, post); err != nil {
			return nil, errors.Wrapf(
				err,
				"userId: %v, lastId: %v, count: %v, postId: %v",
				userId, lastId, count, post.Id,
			)
		}
	}
	return posts, nil
}

func (s *Service) Search(q string, userId uint64, lastId uint64, count int) ([]*models.Post, error) {
	var posts []*models.Post
	r := s.database.Handler().
		Where("(title LIKE ? OR description LIKE ?)", "%"+q+"%", "%"+q+"%").
		Where("id < ? ORDER BY id DESC LIMIT ?", lastId, count).
		Preload("Audio").
		Preload("Image").
		Preload("User").
		Preload("User.Image").
		Find(&posts)
	if r.Error != nil {
		return nil, errors.Wrapf(r.Error, "userId: %v, q: %v, lastId: %v, count: %v", userId, q, lastId, count)
	}
	for _, post := range posts {
		if err := s.attachRelations(userId, post); err != nil {
			return nil, errors.Wrapf(
				err,
				"userId: %v, q: %v, lastId: %v, count: %v, postId: %v",
				userId, q, lastId, count, post.Id,
			)
		}
	}
	return posts, nil
}

func (s *Service) attachRelations(userId uint64, post *models.Post) error {
	if err := s.attachLikes(userId, post); err != nil {
		return errors.WithStack(err)
	}

	return errors.WithStack(s.attachComments(post))
}

func (s *Service) attachLikes(userId uint64, post *models.Post) error {
	var likes []*models.Like

	r := s.database.Handler().
		Where("post_id = ?", post.Id).
		Where("user_id = ?", userId).
		Preload("User").
		Find(&likes)
	if r.Error != nil {
		return errors.Wrapf(r.Error, "userId: %v, postId: %v", userId, post.Id)
	}
	post.Likes = likes

	var count int64
	r = s.database.Handler().Model(&models.Like{}).Where("post_id = ?", post.Id).Count(&count)
	if r.Error != nil {
		return errors.Wrapf(r.Error, "postId: %v", post.Id)
	}
	post.LikesCount = uint64(count)

	return nil
}

func (s *Service) attachComments(post *models.Post) error {
	var count int64
	r := s.database.Handler().Model(&models.Comment{}).Where("post_id = ?", post.Id).Count(&count)
	if r.Error != nil {
		return errors.Wrapf(r.Error, "postId: %v", post.Id)
	}
	post.CommentsCount = uint64(count)

	return nil
}

func (s *Service) Create(post *models.Post) (uint64, error) {
	r := s.database.Handler().Create(post)
	if r.Error != nil {
		return 0, errors.Wrapf(r.Error, "post: %v", post)
	}
	return post.Id, nil
}

func (s *Service) UpdateFields(post *models.Post, title, description string) *models.Post {
	post.Title = title
	post.Description = description
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
		Preload("Audio").
		Preload("Image").
		Preload("User").
		Preload("User.Image").
		First(&post)
	if r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.Wrapf(r.Error, "id: %v", id)
	}
	if err := s.attachRelations(post.UserId, &post); err != nil {
		return nil, errors.Wrapf(err, "userId: %v, postId: %v", post.UserId, post.Id)
	}
	return &post, nil
}

func (s *Service) FindBy(field string, value interface{}) (*models.Post, error) {
	var post models.Post
	r := s.database.Handler().
		Where(field+" = ?", value).
		Preload("Audio").
		Preload("Image").
		Preload("User").
		Preload("User.Image").
		First(&post)
	if r.Error != nil && errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &post, errors.Wrapf(r.Error, "field: %v, value: %v", field, value)
}

func New(database *database.Database) *Service {
	return &Service{database: database}
}
