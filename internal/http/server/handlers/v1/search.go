package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/neatplex/nightel-core/internal/models"
	"github.com/neatplex/nightel-core/internal/services/container"
	"github.com/neatplex/nightel-core/internal/utils"
	"net/http"
)

func Search(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*models.User)

		posts, err := ctr.PostService.Search(
			ctx.QueryParams().Get("q"),
			user.ID,
			utils.StringToID(ctx.QueryParams().Get("lastId"), ^uint64(0)),
			utils.StringToInt(ctx.QueryParams().Get("count"), 10),
		)
		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusOK, map[string][]*models.Post{
			"posts": posts,
		})
	}
}
