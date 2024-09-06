package container

import (
	"github.com/neatplex/nightell-core/internal/config"
	"github.com/neatplex/nightell-core/internal/database"
	"github.com/neatplex/nightell-core/internal/gc"
	"github.com/neatplex/nightell-core/internal/logger"
	"github.com/neatplex/nightell-core/internal/mailer"
	"github.com/neatplex/nightell-core/internal/s3"
	"github.com/neatplex/nightell-core/internal/services/comment"
	"github.com/neatplex/nightell-core/internal/services/file"
	"github.com/neatplex/nightell-core/internal/services/followship"
	"github.com/neatplex/nightell-core/internal/services/like"
	"github.com/neatplex/nightell-core/internal/services/otp"
	"github.com/neatplex/nightell-core/internal/services/post"
	"github.com/neatplex/nightell-core/internal/services/remove"
	"github.com/neatplex/nightell-core/internal/services/setting"
	"github.com/neatplex/nightell-core/internal/services/token"
	"github.com/neatplex/nightell-core/internal/services/user"
)

type Container struct {
	Config            *config.Config
	Logger            *logger.Logger
	S3                *s3.S3
	Database          *database.Database
	Mailer            *mailer.Mailer
	GC                *gc.Gc
	SettingService    *setting.Service
	UserService       *user.Service
	TokenService      *token.Service
	RemoveService     *remove.Service
	PostService       *post.Service
	CommentService    *comment.Service
	FileService       *file.Service
	LikeService       *like.Service
	FollowshipService *followship.Service
	OtpService        *otp.Service
}

func New(
	config *config.Config,
	logger *logger.Logger,
	s3 *s3.S3,
	db *database.Database,
	mailer *mailer.Mailer,
	gc *gc.Gc,
) *Container {
	return &Container{
		Config:            config,
		Logger:            logger,
		S3:                s3,
		Database:          db,
		Mailer:            mailer,
		GC:                gc,
		SettingService:    setting.New(config),
		UserService:       user.New(db, mailer),
		TokenService:      token.New(db),
		RemoveService:     remove.New(db),
		PostService:       post.New(db),
		CommentService:    comment.New(db),
		FileService:       file.New(db, s3),
		LikeService:       like.New(db),
		FollowshipService: followship.New(db),
		OtpService:        otp.New(mailer),
	}
}
