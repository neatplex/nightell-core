package server

import (
	"github.com/labstack/echo/v4/middleware"
	"github.com/neatplex/nightel-core/internal/http/server/handlers"
	"github.com/neatplex/nightel-core/internal/http/server/handlers/v1"
	mw "github.com/neatplex/nightel-core/internal/http/server/middleware"
)

func (s *Server) registerRoutes() {
	s.E.GET("/debug", handlers.Debug)
	s.E.GET("/healthz", handlers.Healthz)

	v1Api := s.E.Group("/api/v1", middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(5)))
	{
		public := v1Api.Group("", middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(1)))
		{
			// auth
			public.POST("/auth/sign-up", v1.AuthSignUp(s.container))
			public.POST("/auth/sign-in", v1.AuthSignIn(s.container))
		}

		private := v1Api.Group("", mw.Authorize(s.container))
		{
			// stories
			private.POST("/stories", v1.StoriesIndex(s.container))
			private.GET("/stories", v1.StoriesIndex(s.container))
			private.POST("/stories", v1.StoriesStore(s.container))
			// files
			private.POST("/files", v1.FilesStore(s.container))
		}
	}
}
