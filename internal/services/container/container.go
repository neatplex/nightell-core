package container

import (
	"github.com/neatplex/nightel-core/internal/database"
	"github.com/neatplex/nightel-core/internal/services/token"
	"github.com/neatplex/nightel-core/internal/services/user"
)

type Container struct {
	UserService  *user.Service
	TokenService *token.Service
}

func New(database *database.Database) *Container {
	return &Container{
		UserService:  user.New(database),
		TokenService: token.New(database),
	}
}
