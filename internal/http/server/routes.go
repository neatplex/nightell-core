package server

import (
	"github.com/labstack/echo/v4/middleware"
	"github.com/neatplex/nightell-core/internal/http/server/handlers"
	"github.com/neatplex/nightell-core/internal/http/server/handlers/v1"
	mw "github.com/neatplex/nightell-core/internal/http/server/middleware"
)

func (s *Server) registerRoutes() {
	s.E.GET("/healthz", handlers.Healthz)

	// delete-account
	s.E.POST("/delete-request", handlers.DeleteRequest(s.container))
	s.E.GET("/delete-account", handlers.DeleteAccount(s.container))

	v1Api := s.E.Group("/api/v1", middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(5)))
	{
		public := v1Api.Group("/", middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(1)))
		{
			// settings
			public.GET("settings", v1.SettingsIndex(s.container))
			// auth
			public.POST("auth/sign-up", v1.AuthSignUp(s.container))
			public.POST("auth/sign-in/email", v1.AuthSignInEmail(s.container))
			public.POST("auth/sign-in/username", v1.AuthSignInUsername(s.container))
			public.POST("auth/sign-in/google", v1.AuthSignInGoogle(s.container))
			public.POST("auth/otp/email/send", v1.AuthOtpEmailSend(s.container))
			public.POST("auth/otp/email/verify", v1.AuthOtpEmailVerify(s.container))
		}

		private := v1Api.Group("/", mw.Authorize(s.container))
		{
			// profile
			private.GET("profile", v1.ProfileShow(s.container))
			private.PATCH("profile/name", v1.ProfileUpdateName(s.container))
			private.PATCH("profile/bio", v1.ProfileUpdateBio(s.container))
			private.PATCH("profile/username", v1.ProfileUpdateUsername(s.container))
			private.PATCH("profile/image", v1.ProfileUpdateImage(s.container))
			private.PATCH("profile/password", v1.ProfileUpdatePassword(s.container))
			private.PATCH("profile/email", v1.ProfileUpdateEmail(s.container))
			private.POST("profile/email/verification", v1.ProfileUpdateEmailVerification(s.container))
			private.DELETE("profile", v1.ProfileDelete(s.container))
			// users
			private.GET("users/:userId", v1.UsersShow(s.container))
			private.GET("users/:userId/followings", v1.UsersFollowings(s.container))
			private.GET("users/:userId/followers", v1.UsersFollowers(s.container))
			private.POST("users/:userId/followers", v1.UsersFollowersStore(s.container))
			private.DELETE("users/:userId/followers", v1.UsersFollowersDelete(s.container))
			// posts
			private.GET("users/:userId/posts", v1.PostsIndex(s.container))
			private.POST("posts", v1.PostsStore(s.container))
			private.GET("posts/:postId", v1.PostsShow(s.container))
			private.PUT("posts/:postId", v1.PostsUpdate(s.container))
			private.DELETE("posts/:postId", v1.PostsDelete(s.container))
			// comments
			private.GET("posts/:postId/comments", v1.CommentsIndexByPost(s.container))
			private.GET("users/:userId/comments", v1.CommentsIndexByUser(s.container))
			private.POST("comments", v1.CommentsStore(s.container))
			private.DELETE("comments/:commentId", v1.CommentsDelete(s.container))
			// likes
			private.GET("posts/:postId/likes", v1.LikesIndex(s.container))
			private.POST("posts/:postId/likes", v1.LikesStoreForPort(s.container))
			private.POST("likes", v1.LikesStore(s.container))
			private.DELETE("likes/:likeId", v1.LikesDelete(s.container))
			// files
			private.POST("files", v1.FilesStore(s.container))
			// search
			private.GET("search", v1.SearchPosts(s.container))
			private.GET("search/posts", v1.SearchPosts(s.container))
			private.GET("search/users", v1.SearchUsers(s.container))
			// feed
			private.GET("feed", v1.Feed(s.container))
		}
	}
}
