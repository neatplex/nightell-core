package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/neatplex/nightel-core/internal/models"
	"github.com/neatplex/nightel-core/internal/services/container"
	"net/http"
)

func StoriesIndex(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*models.User)
		stories, err := ctr.StoryService.Index(user.ID)
		if err != nil {
			return err
		}
		return ctx.JSON(http.StatusCreated, stories)
	}
}
