package app

import (
	"context"
	"github.com/cockroachdb/errors"
	"github.com/neatplex/nightell-core/internal/config"
	"github.com/neatplex/nightell-core/internal/database"
	"github.com/neatplex/nightell-core/internal/gc"
	httpServer "github.com/neatplex/nightell-core/internal/http/server"
	"github.com/neatplex/nightell-core/internal/logger"
	"github.com/neatplex/nightell-core/internal/mailer"
	"github.com/neatplex/nightell-core/internal/s3"
	"github.com/neatplex/nightell-core/internal/services/container"
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
	MySQL      *database.Database
	Container  *container.Container
	Mailer     *mailer.Mailer
	Gc         *gc.Gc
}

// New creates an app from the given configuration file.
func New() (a *App, err error) {
	a = &App{}

	a.Config = config.New()
	if err = a.Config.Init(); err != nil {
		return nil, errors.WithStack(err)
	}
	a.Logger = logger.New(a.Config.Logger.Level, a.Config.Logger.Format, a.Config.Development)
	if err = a.Logger.Init(); err != nil {
		return nil, errors.WithStack(err)
	}
	a.Logger.Debug("app: Config & LoggerProxy initialized")

	a.MySQL = database.New(a.Config, a.Logger)
	a.S3 = s3.New(a.Config, a.Logger)
	a.Container = container.New(a.MySQL, a.S3)
	a.Gc = gc.New(a.MySQL, a.S3, a.Logger)
	a.HttpServer = httpServer.New(a.Config, a.Logger, a.Container)
	a.Mailer = mailer.New(a.Config, a.Logger)
	a.Logger.Debug("app: application modules initialized")
	a.setupSignalListener()

	return a, nil
}

func (a *App) Init() error {
	if err := a.MySQL.Init(); err != nil {
		return errors.WithStack(err)
	}
	if err := a.S3.Init(); err != nil {
		return errors.WithStack(err)
	}
	a.Gc.Init()
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
		a.Logger.Info("app: system call", zap.String("signal", s.String()))
		cancel()
	}()
}

func (a *App) Close() {
	if a.HttpServer != nil {
		a.HttpServer.Close()
	}
	if a.MySQL != nil {
		a.MySQL.Close()
	}
	if a.Logger != nil {
		a.Logger.Close()
	}
}

// Wait avoid dying app and shut it down gracefully on exit signals.
func (a *App) Wait() {
	<-a.context.Done()
}
