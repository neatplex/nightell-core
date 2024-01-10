package database

import (
	"context"
	"database/sql"
	"errors"
	"github.com/neatplex/nightel-core/internal/config"
	"github.com/neatplex/nightel-core/internal/models"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type Database struct {
	handler *gorm.DB
	config  *config.Config
	logger  *zap.Logger
}

func (d *Database) Handler() *gorm.DB {
	return d.handler
}

func (d *Database) Init() {
	timeout := time.Duration(d.config.Database.Timeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	db, err := d.initDatabase(ctx)
	if err != nil {
		d.logger.Fatal("database: cannot connect", zap.Error(err))
	}

	if d.handler, err = gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{}); err != nil {
		d.logger.Fatal("database: cannot initialize gorm", zap.Error(err))
	} else {
		d.logger.Debug("database: gorm connection established successfully")
	}
	d.migrate()
}

func (d *Database) initDatabase(ctx context.Context) (*sql.DB, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, errors.New("database: initial connection timed out")
		default:
			db, err := sql.Open("mysql", d.config.Database.DSN())
			if err != nil {
				return nil, err
			}
			if err = db.Ping(); err == nil {
				return db, nil
			}
			d.logger.Debug("database: trying to connect", zap.Error(err))
			time.Sleep(1 * time.Second)
		}
	}
}

func (d *Database) migrate() {
	err := d.handler.AutoMigrate(
		&models.User{},
		&models.Token{},
		&models.Story{},
		&models.File{},
	)
	if err != nil {
		d.logger.Fatal("database: cannot run migrations", zap.Error(err))
	} else {
		d.logger.Debug("database: migrations ran successfully")
	}
}

func (d *Database) Close() {
	if db, err := d.handler.DB(); err != nil {
		d.logger.Error("database: cannot get DB from GORM to close", zap.Error(err))
	} else {
		if err = db.Close(); err != nil {
			d.logger.Error("database: cannot close database", zap.Error(err))
		}
	}
}

func New(c *config.Config, l *zap.Logger) *Database {
	return &Database{
		config: c,
		logger: l,
	}
}
