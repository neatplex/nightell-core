package server

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/neatplex/nightel-core/internal/config"
	"github.com/neatplex/nightel-core/internal/http/server/validator"
	"github.com/neatplex/nightel-core/internal/logger"
	"github.com/neatplex/nightel-core/internal/services/container"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Server struct {
	E         *echo.Echo
	config    *config.Config
	l         *logger.Logger
	container *container.Container
}

func New(config *config.Config, log *logger.Logger, container *container.Container) *Server {
	e := echo.New()

	e.HideBanner = true
	e.Server.ReadTimeout, _ = time.ParseDuration(config.HTTPServer.ReadTimeout)
	e.Server.WriteTimeout, _ = time.ParseDuration(config.HTTPServer.WriteTimeout)
	e.Server.ReadHeaderTimeout, _ = time.ParseDuration(config.HTTPServer.ReadHeaderTimeout)
	e.Server.IdleTimeout, _ = time.ParseDuration(config.HTTPServer.IdleTimeout)
	e.Validator = validator.New()

	return &Server{E: e, config: config, l: log, container: container}
}

func (s *Server) Serve() {
	s.E.Use(middleware.CORS())
	s.E.Use(middleware.Logger())
	s.E.Use(middleware.Static("web"))
	s.E.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))
	s.E.Use(middleware.BodyLimit("20M"))

	s.registerRoutes()

	go func() {
		listen := s.config.HTTPServer.Listen
		if err := s.E.Start(listen); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.l.Fatal("cannot start the server", zap.String("listen", listen), zap.Error(err))
		}
	}()
}

func (s *Server) Close() {
	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.E.Shutdown(c); err != nil {
		s.l.Warn("cannot close the http server", zap.Error(err))
	}
	s.l.Debug("http server closed successfully")
}
