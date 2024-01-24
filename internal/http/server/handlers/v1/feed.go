package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/neatplex/nightel-core/internal/models"
	"github.com/neatplex/nightel-core/internal/services/container"
	"github.com/neatplex/nightel-core/internal/utils"
	"net/http"
)

func Feed(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*models.User)

		lastId := ^uint64(0)
		requestLastId := ctx.QueryParams().Get("lastId")
		if requestLastId != "" {
			story, err := ctr.StoryService.FindById(utils.StringToID(requestLastId))
			if err != nil {
				return err
			}
			if story != nil {
				lastId = story.ID
			}
		}

		stories, err := ctr.StoryService.Feed(user.ID, lastId)
		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusOK, struct {
			Stories []*models.Story `json:"stories"`
		}{
			Stories: stories,
		})
	}
}
