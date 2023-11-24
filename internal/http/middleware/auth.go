package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/neatplex/nightel-core/internal/services/container"
)

func Authorize(ctr *container.Container) func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			authHeader := ctx.Request().Header.Get("Authorization")
			if len(authHeader) < len("Bearer X") {
				return echo.ErrUnauthorized
			}

			token, err := ctr.TokenService.FindByValue(authHeader[len("Bearer "):])
			if err != nil {
				return err
			}
			if token == nil {
				return echo.ErrUnauthorized
			}

			ctx.Set("user", token.User)

			return next(ctx)
		}
	}
}
