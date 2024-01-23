package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/neatplex/nightel-core/internal/services/container"
	"net/http"
)

func Feed(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		lastId := ^uint64(0)
		identity := ctx.QueryParams().Get("lastIdentity")
		if identity != "" {
			story, err := ctr.StoryService.FindByIdentity(identity)
			if err != nil {
				return err
			}
			if story != nil {
				lastId = story.ID
			}
		}

		stories, err := ctr.StoryService.Feed(lastId)
		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusOK, stories)
	}
}
