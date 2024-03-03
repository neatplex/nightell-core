package server

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/neatplex/nightel-core/internal/config"
	middleware2 "github.com/neatplex/nightel-core/internal/http/server/middleware"
	"github.com/neatplex/nightel-core/internal/http/server/validator"
	"github.com/neatplex/nightel-core/internal/logger"
	"github.com/neatplex/nightel-core/internal/services/container"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	E         *echo.Echo
	config    *config.Config
	l         *logger.Logger
	container *container.Container
}

func New(config *config.Config, logger *logger.Logger, container *container.Container) *Server {
	e := echo.New()
	e.HideBanner = true
	e.Server.ReadTimeout = time.Duration(config.HttpServer.ReadTimeout) * time.Second
	e.Server.WriteTimeout = time.Duration(config.HttpServer.WriteTimeout) * time.Second
	e.Server.ReadHeaderTimeout = time.Duration(config.HttpServer.ReadHeaderTimeout) * time.Second
	e.Server.IdleTimeout = time.Duration(config.HttpServer.IdleTimeout) * time.Second
	e.Validator = validator.New()
	return &Server{E: e, config: config, l: logger, container: container}
}

func (s *Server) Serve() {
	s.E.Use(middleware2.Logger(s.l))
	s.E.Use(middleware.CORS())
	s.E.Use(middleware.Static("web"))
	s.E.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))
	s.E.Use(middleware.BodyLimit("20M"))

	s.registerRoutes()

	go func() {
		listen := s.config.HttpServer.Host + ":" + strconv.Itoa(s.config.HttpServer.Port)
		if err := s.E.Start(listen); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.l.Error("cannot start the http server", zap.String("listen", listen), zap.Error(err))
		}
	}()
}

func (s *Server) Close() {
	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.E.Shutdown(c); err != nil {
		s.l.Error("cannot close the http server", zap.Error(err))
	}
	s.l.Debug("http server closed successfully")
}
