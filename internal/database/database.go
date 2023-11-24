package database

import (
	"database/sql"
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

func (d *Database) Connect() {
	var (
		db  *sql.DB
		err error
	)

	for {
		db, err = sql.Open("mysql", d.config.Database.DSN())
		if err != nil {
			d.logger.Fatal("cannot init connection to mysql", zap.Error(err))
		}

		if err = db.Ping(); err == nil {
			break
		}

		time.Sleep(1 * time.Second)
	}

	if d.handler, err = gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{}); err != nil {
		d.logger.Fatal("cannot connect to database", zap.Error(err))
	}
}

func (d *Database) Migrate() {
	err := d.handler.AutoMigrate(
		&models.User{},
		&models.Token{},
	)
	if err != nil {
		d.logger.Fatal("cannot run database migrations", zap.Error(err))
	}
}

func (d *Database) Close() {
	if db, err := d.handler.DB(); err != nil {
		d.logger.Error("cannot get DB from GORM to close", zap.Error(err))
	} else {
		if err = db.Close(); err != nil {
			d.logger.Error("cannot close database", zap.Error(err))
		}
	}
}

func New(c *config.Config, l *zap.Logger) *Database {
	return &Database{
		config: c,
		logger: l,
	}
}
