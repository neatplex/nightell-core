package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/neatplex/nightel-core/internal/services/container"
	"net/http"
)

func StoriesIndex(_ *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return ctx.JSON(http.StatusCreated, map[string]interface{}{
			"message":          "soon!",
			"authenticated-as": ctx.Get("user"),
		})
	}
}
