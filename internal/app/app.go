package app

import (
	"context"
	"github.com/cockroachdb/errors"
	"github.com/neatplex/nightell-core/internal/config"
	"github.com/neatplex/nightell-core/internal/container"
	"github.com/neatplex/nightell-core/internal/database"
	"github.com/neatplex/nightell-core/internal/gc"
	httpServer "github.com/neatplex/nightell-core/internal/http/server"
	"github.com/neatplex/nightell-core/internal/logger"
	"github.com/neatplex/nightell-core/internal/mailer"
	"github.com/neatplex/nightell-core/internal/s3"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

// App integrates the modules to serve.
type App struct {
	context    context.Context
	HttpServer *httpServer.Server
	Container  *container.Container
}

// New creates an app from the given configuration file.
func New() (a *App, err error) {
	a = &App{}

	c := config.New()
	if err = c.Init(); err != nil {
		return a, errors.WithStack(err)
	}

	l := logger.New(c.Logger.Level, c.Logger.Format, c.Development)
	if err = l.Init(); err != nil {
		return a, errors.WithStack(err)
	}

	db := database.New(c, l)
	if err = db.Init(); err != nil {
		return a, errors.WithStack(err)
	}

	awsS3 := s3.New(c, l)
	if err = awsS3.Init(); err != nil {
		return a, errors.WithStack(err)
	}

	s3gc := gc.New(db, awsS3, l)
	s3gc.Init()

	m := mailer.New(c, l)

	a.Container = container.New(c, l, awsS3, db, m, s3gc)
	a.HttpServer = httpServer.New(a.Container)

	a.setupSignalListener()

	l.Debug("app: modules initialized")

	return a, nil
}

func (a *App) Init() error {

	return nil
}

// setupSignalListener sets up a listener to exit signals from os and closes the app gracefully.
func (a *App) setupSignalListener() {
	var cancel context.CancelFunc
	a.context, cancel = context.WithCancel(context.Background())

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	go func() {
		s := <-signalChannel
		a.Container.Logger.Info("app: system call", zap.String("signal", s.String()))
		cancel()
	}()
}

func (a *App) Close() {
	if a.HttpServer != nil {
		a.HttpServer.Close()
	}
	if a.Container.DB != nil {
		a.Container.DB.Close()
	}
	if a.Container.Logger != nil {
		a.Container.Logger.Close()
	}
}

// Wait avoid dying app and shut it down gracefully on exit signals.
func (a *App) Wait() {
	<-a.context.Done()
}
