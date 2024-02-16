package followship

import (
	"fmt"
	"github.com/neatplex/nightel-core/internal/database"
	"github.com/neatplex/nightel-core/internal/models"
)

type Service struct {
	database *database.Database
}

func (s *Service) FindByIds(followerID, followeeID uint64) (*models.Followship, error) {
	var followship models.Followship
	r := s.database.Handler().
		Where("follower_id = ?", followerID).
		Where("followee_id = ?", followeeID).
		Find(&followship)
	if r.Error != nil {
		return nil, fmt.Errorf("services: followship: FindByIds: %s", r.Error)
	}
	return &followship, nil
}

func (s *Service) IndexFollowers(userId uint64, lastId uint64, count int) ([]*models.Followship, error) {
	var followships []*models.Followship
	r := s.database.Handler().
		Preload("Follower").
		Where("followee_id = ?", userId).
		Where("id < ? ORDER BY id DESC LIMIT ?", lastId, count).
		Find(&followships)
	if r.Error != nil {
		return nil, fmt.Errorf("services: followship: IndexFollowers: %s", r.Error)
	}
	return followships, nil
}

func (s *Service) IndexFollowings(userId uint64, lastId uint64, count int) ([]*models.Followship, error) {
	var followships []*models.Followship
	r := s.database.Handler().
		Preload("Followee").
		Where("follower_id = ?", userId).
		Where("id < ? ORDER BY id DESC LIMIT ?", lastId, count).
		Find(&followships)
	if r.Error != nil {
		return nil, fmt.Errorf("services: followship: IndexFollowings: %s", r.Error)
	}
	return followships, nil
}

func (s *Service) CountFollowers(userId uint64) (int64, error) {
	var count int64
	r := s.database.Handler().
		Model(&models.Followship{}).
		Where("followee_id = ?", userId).
		Count(&count)
	if r.Error != nil {
		return -1, fmt.Errorf("services: followship: CountFollowers: %s", r.Error)
	}
	return count, nil
}

func (s *Service) CountFollowings(userId uint64) (int64, error) {
	var count int64
	r := s.database.Handler().
		Model(&models.Followship{}).
		Where("follower_id = ?", userId).
		Count(&count)
	if r.Error != nil {
		return -1, fmt.Errorf("services: followship: CountFollowings: %s", r.Error)
	}
	return count, nil
}

func (s *Service) Create(followeeID, followerID uint64) (*models.Followship, error) {
	var followship models.Followship
	r := s.database.Handler().Preload("Followee").FirstOrCreate(&followship, &models.Followship{
		FollowerID: followerID,
		FolloweeID: followeeID,
	})
	if r.Error != nil {
		return nil, fmt.Errorf("services: followship: Create: %v, err: %v", followship, r.Error)
	}
	return &followship, r.Error
}

func (s *Service) Delete(id uint64) error {
	r := s.database.Handler().Delete(&models.Followship{}, id)
	if r.Error != nil {
		return fmt.Errorf("services: followship: Delete #%v, err: %v", id, r.Error)
	}
	return r.Error
}

func New(database *database.Database) *Service {
	return &Service{database: database}
}
