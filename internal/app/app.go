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
	Config     *config.Config
	Logger     *logger.Logger
	S3         *s3.S3
	HttpServer *httpServer.Server
	Database   *database.Database
	Container  *container.Container
}

// New creates an app from the given configuration file.
func New(configPath string) (a *App, err error) {
	a = &App{}

	a.Config, err = config.New(configPath)
	if err != nil {
		return nil, err
	}
	a.Logger, err = logger.New(a.Config, a.ShutdownModules)
	if err != nil {
		return nil, err
	}
	a.Logger.Debug("app: config & logger initialized")

	a.Database = database.New(a.Config, a.Logger)
	a.S3 = s3.New(a.Config, a.Logger)
	a.Container = container.New(a.Database, a.S3)
	a.HttpServer = httpServer.New(a.Config, a.Logger, a.Container)

	a.Logger.Debug("app: application modules initialized")

	a.setupSignalListener()

	return a, nil
}

// Boot makes sure the critical modules and external sources work fine.
func (a *App) Boot() {
	a.Database.Init()
	a.S3.Init()
	a.HttpServer.Serve()
}

// setupSignalListener sets up a listener to exit signals from os and closes the app gracefully.
func (a *App) setupSignalListener() {
	var cancel context.CancelFunc
	a.context, cancel = context.WithCancel(context.Background())

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	go func() {
		s := <-signalChannel
		a.Logger.Info("app: system call", zap.String("signal", s.String()))
		cancel()
	}()
}

func (a *App) ShutdownModules() {
	if a.HttpServer != nil {
		a.HttpServer.Close()
	}
	if a.Database != nil {
		a.Database.Close()
	}
}

// Wait avoid dying app and shut it down gracefully on exit signals.
func (a *App) Wait() {
	<-a.context.Done()

	a.ShutdownModules()

	if a.Logger != nil {
		a.Logger.Close()
	}
}
