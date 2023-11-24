package server

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/neatplex/nightel-core/internal/config"
	"github.com/neatplex/nightel-core/internal/http/handlers"
	"github.com/neatplex/nightel-core/internal/http/handlers/v1"
	mw "github.com/neatplex/nightel-core/internal/http/middleware"
	"github.com/neatplex/nightel-core/internal/http/validator"
	"github.com/neatplex/nightel-core/internal/services/container"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Server struct {
	E         *echo.Echo
	config    *config.Config
	log       *zap.Logger
	container *container.Container
}

func New(config *config.Config, log *zap.Logger, container *container.Container) *Server {
	e := echo.New()

	e.HideBanner = true
	e.Server.ReadTimeout, _ = time.ParseDuration(config.HTTPServer.ReadTimeout)
	e.Server.WriteTimeout, _ = time.ParseDuration(config.HTTPServer.WriteTimeout)
	e.Server.ReadHeaderTimeout, _ = time.ParseDuration(config.HTTPServer.ReadHeaderTimeout)
	e.Server.IdleTimeout, _ = time.ParseDuration(config.HTTPServer.IdleTimeout)
	e.Validator = validator.New()

	return &Server{E: e, config: config, log: log, container: container}
}

func (s *Server) Serve() {
	s.E.Use(middleware.CORS())
	s.E.Use(middleware.Logger())
	s.E.Use(middleware.Static("web"))
	s.E.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))
	s.E.Use(middleware.BodyLimit("20M"))

	s.E.GET("/healthz", handlers.Healthz)

	v1Api := s.E.Group("/api/v1", middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(5)))
	{
		public := v1Api.Group("", middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(1)))
		{
			public.POST("/auth/sign-up", v1.AuthSignUp(s.container))
			public.POST("/auth/sign-in", v1.AuthSignIn(s.container))
		}

		private := v1Api.Group("", mw.Authorize(s.container))
		{
			private.GET("/stories", v1.StoriesIndex(s.container))
		}
	}

	go func() {
		listen := s.config.HTTPServer.Listen
		if err := s.E.Start(listen); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Fatal("cannot start the server", zap.String("listen", listen), zap.Error(err))
		}
	}()
}

func (s *Server) Close() {
	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.E.Shutdown(c); err != nil {
		s.log.Warn("cannot close the http server", zap.Error(err))
	}
	s.log.Debug("http server closed successfully")
}
