package container

import (
	"github.com/neatplex/nightell-core/internal/database"
	"github.com/neatplex/nightell-core/internal/mailer"
	"github.com/neatplex/nightell-core/internal/s3"
	"github.com/neatplex/nightell-core/internal/services/chat"
	"github.com/neatplex/nightell-core/internal/services/file"
	"github.com/neatplex/nightell-core/internal/services/followship"
	"github.com/neatplex/nightell-core/internal/services/like"
	"github.com/neatplex/nightell-core/internal/services/message"
	"github.com/neatplex/nightell-core/internal/services/post"
	"github.com/neatplex/nightell-core/internal/services/remove"
	"github.com/neatplex/nightell-core/internal/services/token"
	"github.com/neatplex/nightell-core/internal/services/user"
)

type Container struct {
	UserService       *user.Service
	TokenService      *token.Service
	RemoveService     *remove.Service
	PostService       *post.Service
	FileService       *file.Service
	LikeService       *like.Service
	FollowshipService *followship.Service
	MessageService    *message.Service
	ChatService       *chat.Service
}

func New(d *database.Database, s3 *s3.S3, m *mailer.Mailer) *Container {
	messageService := message.New(d)
	return &Container{
		UserService:       user.New(d, m),
		TokenService:      token.New(d),
		RemoveService:     remove.New(d),
		PostService:       post.New(d),
		FileService:       file.New(d, s3),
		LikeService:       like.New(d),
		FollowshipService: followship.New(d),
		MessageService:    messageService,
		ChatService:       chat.New(d, messageService),
	}
}
