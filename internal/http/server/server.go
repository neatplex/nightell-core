package server

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/neatplex/nightell-core/internal/container"
	middleware2 "github.com/neatplex/nightell-core/internal/http/server/middleware"
	"github.com/neatplex/nightell-core/internal/http/server/validator"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	E         *echo.Echo
	container *container.Container
}

func New(c *container.Container) *Server {
	e := echo.New()
	e.HideBanner = true
	e.Server.ReadTimeout = time.Duration(c.Config.HttpServer.ReadTimeout) * time.Second
	e.Server.WriteTimeout = time.Duration(c.Config.HttpServer.WriteTimeout) * time.Second
	e.Server.ReadHeaderTimeout = time.Duration(c.Config.HttpServer.ReadHeaderTimeout) * time.Second
	e.Server.IdleTimeout = time.Duration(c.Config.HttpServer.IdleTimeout) * time.Second
	e.Validator = validator.New()
	return &Server{E: e, container: c}
}

func (s *Server) Serve() {
	l := s.container.Logger
	c := s.container.Config

	s.E.Use(middleware2.Logger(l))
	s.E.Use(middleware.CORS())
	s.E.Use(middleware.Static("web"))
	s.E.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))
	s.E.Use(middleware.BodyLimit("6M"))

	s.registerRoutes()

	go func() {
		listen := c.HttpServer.Host + ":" + strconv.Itoa(c.HttpServer.Port)
		if err := s.E.Start(listen); err != nil && !errors.Is(err, http.ErrServerClosed) {
			l.Error("cannot start the http server", zap.String("listen", listen), zap.Error(err))
		}
	}()
}

func (s *Server) Close() {
	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.E.Shutdown(c); err != nil {
		s.container.Logger.Error("cannot close the http server", zap.Error(err))
	}
	s.container.Logger.Debug("http server closed successfully")
}
