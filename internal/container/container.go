package container

import (
	"github.com/neatplex/nightel-core/internal/database"
	"github.com/neatplex/nightel-core/internal/s3"
	"github.com/neatplex/nightel-core/internal/services/file"
	"github.com/neatplex/nightel-core/internal/services/story"
	"github.com/neatplex/nightel-core/internal/services/token"
	"github.com/neatplex/nightel-core/internal/services/user"
)

type Container struct {
	UserService  *user.Service
	TokenService *token.Service
	StoryService *story.Service
	FileService  *file.Service
}

func New(database *database.Database, s3 *s3.S3) *Container {
	return &Container{
		UserService:  user.New(database),
		TokenService: token.New(database),
		StoryService: story.New(database),
		FileService:  file.New(database, s3),
	}
}
