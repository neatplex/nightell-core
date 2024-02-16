package server

import (
	"github.com/labstack/echo/v4/middleware"
	"github.com/neatplex/nightel-core/internal/http/server/handlers"
	"github.com/neatplex/nightel-core/internal/http/server/handlers/v1"
	mw "github.com/neatplex/nightel-core/internal/http/server/middleware"
)

func (s *Server) registerRoutes() {
	s.E.GET("/healthz", handlers.Healthz)

	v1Api := s.E.Group("/api/v1", middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(5)))
	{
		public := v1Api.Group("/", middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(1)))
		{
			// auth
			public.POST("auth/sign-up", v1.AuthSignUp(s.container))
			public.POST("auth/sign-in/email", v1.AuthSignInEmail(s.container))
			public.POST("auth/sign-in/username", v1.AuthSignInUsername(s.container))
		}

		private := v1Api.Group("/", mw.Authorize(s.container))
		{
			// profile
			private.GET("profile", v1.ProfileShow())
			private.PATCH("profile/name", v1.ProfileUpdateName(s.container))
			private.PATCH("profile/bio", v1.ProfileUpdateBio(s.container))
			private.PATCH("profile/username", v1.ProfileUpdateUsername(s.container))
			// users
			private.GET("users/:userId", v1.UsersShow(s.container))
			// posts
			private.GET("users/:userId/posts", v1.PostsIndex(s.container))
			private.POST("posts", v1.PostsStore(s.container))
			private.PATCH("posts/:postId/caption", v1.PostsUpdateCaption(s.container))
			private.DELETE("posts/:postId", v1.PostsDelete(s.container))
			// likes
			private.GET("posts/:postId/likes", v1.LikesIndex(s.container))
			private.POST("posts/:postId/likes", v1.LikesStore(s.container))
			private.DELETE("likes/:likeId", v1.LikesDelete(s.container))
			// files
			private.POST("files", v1.FilesStore(s.container, s.l))
			// feed
			private.GET("feed", v1.Feed(s.container))
		}
	}
}
