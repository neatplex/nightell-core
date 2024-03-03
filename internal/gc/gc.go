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
		var post *models.Post
		r = g.database.Handler().
			Where("audio_id = ?", file.ID).
			Or("image_id = ?", file.ID).
			First(&post)
		if r.Error != nil && !errors.Is(r.Error, gorm.ErrRecordNotFound) {
			g.l.Error("gc: cannot fetch post", zap.Error(errors.WithStack(r.Error)))
			return
		}
		if post == nil || post.ID == 0 {
			if err := g.s3.Delete(file.Path); err != nil {
				g.l.Error("gc: cannot delete s3 file", zap.Error(errors.WithStack(err)))
				return
			}
			r = g.database.Handler().Delete(&file)
			if r.Error != nil {
				g.l.Error("gc: cannot delete db file", zap.Error(errors.WithStack(r.Error)))
			} else {
				g.l.Info("gc: cleaned", zap.Uint64("id", file.ID), zap.String("path", file.Path))
			}
		}
	}
}

func New(database *database.Database, s3 *s3.S3, l *logger.Logger) *Gc {
	return &Gc{database: database, s3: s3, l: l}
}
