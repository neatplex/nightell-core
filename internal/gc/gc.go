package gc

import (
	"github.com/cockroachdb/errors"
	"github.com/neatplex/nightell-core/internal/database"
	"github.com/neatplex/nightell-core/internal/logger"
	"github.com/neatplex/nightell-core/internal/models"
	"github.com/neatplex/nightell-core/internal/s3"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type Gc struct {
	database *database.Database
	s3       *s3.S3
	l        *logger.Logger
}

func (g *Gc) Init() {
	g.cleanup()
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			g.cleanup()
		}
	}()
}

func (g *Gc) cleanup() {
	var files []*models.File
	r := g.database.Handler().Find(&files)
	if r.Error != nil {
		g.l.Error("gc: cannot fetch files", zap.Error(errors.WithStack(r.Error)))
		return
	}

	for _, file := range files {
		if g.IsUsedForUsers(file) {
			continue
		}
		if g.IsUsedForPosts(file) {
			continue
		}

		if err := g.s3.Delete(file.Path); err != nil {
			g.l.Error("gc: cannot delete s3 file", zap.Error(errors.WithStack(err)))
			return
		}

		r = g.database.Handler().Delete(&file)
		if r.Error != nil {
			g.l.Error("gc: cannot delete db file", zap.Error(errors.WithStack(r.Error)))
		} else {
			g.l.Info("gc: cleaned", zap.Uint64("id", file.Id), zap.String("path", file.Path))
		}
	}
}

func (g *Gc) IsUsedForPosts(file *models.File) bool {
	var post models.Post
	r := g.database.Handler().
		Where("audio_id = ?", file.Id).
		Or("image_id = ?", file.Id).
		First(&post)
	if r.Error != nil && !errors.Is(r.Error, gorm.ErrRecordNotFound) {
		g.l.Error("gc: query failed", zap.Error(errors.WithStack(r.Error)))
	}
	return post.Id != 0
}

func (g *Gc) IsUsedForUsers(file *models.File) bool {
	var user models.User
	r := g.database.Handler().Where("image_id = ?", file.Id).First(&user)
	if r.Error != nil && !errors.Is(r.Error, gorm.ErrRecordNotFound) {
		g.l.Error("gc: query failed", zap.Error(errors.WithStack(r.Error)))
	}
	return user.Id != 0
}

func New(database *database.Database, s3 *s3.S3, l *logger.Logger) *Gc {
	return &Gc{database: database, s3: s3, l: l}
}
