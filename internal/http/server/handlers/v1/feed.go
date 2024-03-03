package v1

import (
	"github.com/cockroachdb/errors"
	"github.com/labstack/echo/v4"
	"github.com/neatplex/nightell-core/internal/models"
	"github.com/neatplex/nightell-core/internal/services/container"
	"github.com/neatplex/nightell-core/internal/utils"
	"net/http"
)

func Feed(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*models.User)

		posts, err := ctr.PostService.Feed(
			user.ID,
			utils.StringToID(ctx.QueryParams().Get("lastId"), ^uint64(0)),
			utils.StringToInt(ctx.QueryParams().Get("count"), 10),
		)
		if err != nil {
			return errors.WithStack(err)
		}

		return ctx.JSON(http.StatusOK, map[string][]*models.Post{
			"posts": posts,
		})
	}
}
