package followship

import (
	"github.com/cockroachdb/errors"
	"github.com/neatplex/nightell-core/internal/database"
	"github.com/neatplex/nightell-core/internal/models"
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
	return &followship, errors.Wrapf(r.Error, "followerId: %v, followeeId: %v", followerID, followeeID)
}

func (s *Service) IndexFollowers(userId uint64, lastUserId uint64, count int) ([]*models.Followship, error) {
	var followships []*models.Followship
	r := s.database.Handler().
		Preload("Follower").
		Where("followee_id = ?", userId).
		Where("follower_id < ? ORDER BY follower_id DESC LIMIT ?", lastUserId, count).
		Find(&followships)
	return followships, errors.Wrapf(r.Error, "userId: %v, lastUserId: %v, count: %v", userId, lastUserId, count)
}

func (s *Service) IndexFollowings(userId uint64, lastUserId uint64, count int) ([]*models.Followship, error) {
	var followships []*models.Followship
	r := s.database.Handler().
		Preload("Followee").
		Where("follower_id = ?", userId).
		Where("followee_id < ? ORDER BY followee_id DESC LIMIT ?", lastUserId, count).
		Find(&followships)
	return followships, errors.Wrapf(r.Error, "userId: %v, lastUserId: %v, count: %v", userId, lastUserId, count)
}

func (s *Service) CountFollowers(userId uint64) (int64, error) {
	var count int64
	r := s.database.Handler().
		Model(&models.Followship{}).
		Where("followee_id = ?", userId).
		Count(&count)
	return count, errors.Wrapf(r.Error, "userId: %v", userId)
}

func (s *Service) CountFollowings(userId uint64) (int64, error) {
	var count int64
	r := s.database.Handler().
		Model(&models.Followship{}).
		Where("follower_id = ?", userId).
		Count(&count)
	return count, errors.Wrapf(r.Error, "userId: %v", userId)
}

func (s *Service) Create(followeeID, followerID uint64) (*models.Followship, error) {
	var followship models.Followship
	r := s.database.Handler().Preload("Followee").FirstOrCreate(&followship, &models.Followship{
		FollowerID: followerID,
		FolloweeID: followeeID,
	})
	return &followship, errors.Wrapf(r.Error, "followerId: %v, followeeId: %v", followerID, followeeID)
}

func (s *Service) Delete(id uint64) error {
	r := s.database.Handler().Delete(&models.Followship{}, id)
	return errors.Wrapf(r.Error, "id: %v", id)
}

func New(database *database.Database) *Service {
	return &Service{database: database}
}
