package app

import (
	"context"
	"github.com/neatplex/nightel-core/internal/config"
	"github.com/neatplex/nightel-core/internal/database"
	httpServer "github.com/neatplex/nightel-core/internal/http/server"
	"github.com/neatplex/nightel-core/internal/logger"
	"github.com/neatplex/nightel-core/internal/s3"
	"github.com/neatplex/nightel-core/internal/services/container"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

// App integrates the modules to serve.
type App struct {
	context    context.Context
	config     *config.Config
	logger     *logger.Logger
	S3         *s3.S3
	httpServer *httpServer.Server
	database   *database.Database
	container  *container.Container
}

// New creates an app from the given configuration file.
func New(configPath string) (a *App, err error) {
	a = &App{}

	a.config, err = config.New(configPath)
	if err != nil {
		return nil, err
	}
	a.logger, err = logger.New(a.config, a.ShutdownModules)
	if err != nil {
		return nil, err
	}
	a.logger.Debug("app: config & logger initialized")

	a.database = database.New(a.config, a.logger)
	a.S3 = s3.New(a.config, a.logger)
	a.container = container.New(a.database, a.S3)
	a.httpServer = httpServer.New(a.config, a.logger, a.container)

	a.logger.Debug("app: application modules initialized")

	a.setupSignalListener()

	return a, nil
}

// Boot makes sure the critical modules and external sources work fine.
func (a *App) Boot() {
	a.database.Init()
	a.S3.Init()
	a.httpServer.Serve()
}

// setupSignalListener sets up a listener to exit signals from os and closes the app gracefully.
func (a *App) setupSignalListener() {
	var cancel context.CancelFunc
	a.context, cancel = context.WithCancel(context.Background())

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	go func() {
		s := <-signalChannel
		a.logger.Info("app: system call", zap.String("signal", s.String()))
		cancel()
	}()
}

func (a *App) ShutdownModules() {
	if a.httpServer != nil {
		a.httpServer.Close()
	}
	if a.database != nil {
		a.database.Close()
	}
}

// Wait avoid dying app and shut it down gracefully on exit signals.
func (a *App) Wait() {
	<-a.context.Done()

	a.ShutdownModules()

	if a.logger != nil {
		a.logger.Close()
	}
}
