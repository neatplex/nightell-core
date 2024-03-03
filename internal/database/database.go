package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/neatplex/nightell-core/internal/config"
	"github.com/neatplex/nightell-core/internal/logger"
	"github.com/neatplex/nightell-core/internal/models"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type Database struct {
	handler *gorm.DB
	config  *config.Config
	l       *logger.Logger
}

func (d *Database) Handler() *gorm.DB {
	return d.handler
}

func (d *Database) Init() error {
	timeout := time.Duration(d.config.MySQL.Timeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	db, err := d.initDatabase(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	if d.handler, err = gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{}); err != nil {
		return errors.WithStack(err)
	} else {
		d.l.Debug("gorm connection established successfully")
	}
	return errors.WithStack(d.migrate())
}

func (d *Database) initDatabase(ctx context.Context) (db *sql.DB, err error) {
	for {
		select {
		case <-ctx.Done():
			return nil, errors.New("initial connection timed out")
		default:
			if db == nil {
				db, err = sql.Open("mysql", d.dsn())
				if err != nil {
					return nil, errors.Wrapf(err, "cannot open connection to %s", d.dsn())
				}
			}
			if err = db.Ping(); err == nil {
				return db, nil
			}
			d.l.Debug("database: trying to connect", zap.Error(err))
			time.Sleep(1 * time.Second)
		}
	}
}

func (d *Database) dsn() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true&multiStatements=true&interpolateParams=true&collation=%s",
		d.config.MySQL.User,
		d.config.MySQL.Password,
		d.config.MySQL.Host,
		d.config.MySQL.Port,
		d.config.MySQL.Name,
		"utf8mb4_general_ci",
	)
}

func (d *Database) migrate() error {
	err := d.handler.AutoMigrate(
		&models.User{},
		&models.Token{},
		&models.Post{},
		&models.File{},
		&models.Like{},
		&models.Followship{},
	)
	if err != nil {
		return errors.WithStack(err)
	}
	d.l.Debug("database: migrations ran successfully")
	return nil
}

func (d *Database) Close() {
	if db, err := d.handler.DB(); err != nil {
		d.l.Error("database: cannot get DB from GORM to close", zap.Error(err))
	} else {
		if err = db.Close(); err != nil {
			d.l.Error("database: cannot close database", zap.Error(err))
		}
	}
}

func New(c *config.Config, l *logger.Logger) *Database {
	return &Database{config: c, l: l}
}
