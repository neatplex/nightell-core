package v1

import (
	"github.com/cockroachdb/errors"
	"github.com/labstack/echo/v4"
	"github.com/neatplex/nightell-core/internal/container"
	"github.com/neatplex/nightell-core/internal/models"
	"github.com/neatplex/nightell-core/internal/utils"
	"net/http"
)

func SearchPosts(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*models.User)

		posts, err := ctr.PostService.Search(
			ctx.QueryParams().Get("q"),
			user.ID,
			utils.StringToID(ctx.QueryParams().Get("lastId"), ^uint64(0)),
			utils.StringToInt(ctx.QueryParams().Get("count"), 100, 10),
		)
		if err != nil {
			return errors.WithStack(err)
		}

		return ctx.JSON(http.StatusOK, map[string][]*models.Post{
			"posts": posts,
		})
	}
}

func SearchUsers(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		users, err := ctr.UserService.Search(
			ctx.QueryParams().Get("q"),
			utils.StringToID(ctx.QueryParams().Get("lastId"), ^uint64(0)),
			utils.StringToInt(ctx.QueryParams().Get("count"), 100, 10),
		)
		if err != nil {
			return errors.WithStack(err)
		}
		return ctx.JSON(http.StatusOK, map[string][]*models.User{
			"users": users,
		})
	}
}
