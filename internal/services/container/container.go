package container

import (
	"github.com/neatplex/nightell-core/internal/database"
	"github.com/neatplex/nightell-core/internal/s3"
	"github.com/neatplex/nightell-core/internal/services/file"
	"github.com/neatplex/nightell-core/internal/services/followship"
	"github.com/neatplex/nightell-core/internal/services/like"
	"github.com/neatplex/nightell-core/internal/services/post"
	"github.com/neatplex/nightell-core/internal/services/token"
	"github.com/neatplex/nightell-core/internal/services/user"
)

type Container struct {
	UserService       *user.Service
	TokenService      *token.Service
	PostService       *post.Service
	FileService       *file.Service
	LikeService       *like.Service
	FollowshipService *followship.Service
}

func New(d *database.Database, s3 *s3.S3) *Container {
	return &Container{
		UserService:       user.New(d),
		TokenService:      token.New(d),
		PostService:       post.New(d),
		FileService:       file.New(d, s3),
		LikeService:       like.New(d),
		FollowshipService: followship.New(d),
	}
}
