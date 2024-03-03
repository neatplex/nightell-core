package v1

import (
	"github.com/cockroachdb/errors"
	"github.com/labstack/echo/v4"
	"github.com/neatplex/nightel-core/internal/models"
	"github.com/neatplex/nightel-core/internal/services/container"
	"github.com/neatplex/nightel-core/internal/utils"
	"net/http"
)

func LikesIndex(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		posts, err := ctr.LikeService.IndexByPostIDWithUser(
			utils.StringToID(ctx.Param("postId"), 0),
			utils.StringToID(ctx.QueryParams().Get("lastId"), ^uint64(0)),
			utils.StringToInt(ctx.QueryParams().Get("count"), 10),
		)
		if err != nil {
			return errors.WithStack(err)
		}
		return ctx.JSON(http.StatusCreated, map[string]interface{}{
			"likes": posts,
		})
	}
}

func LikesStore(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*models.User)

		post, err := ctr.PostService.FindById(utils.StringToID(ctx.Param("postId"), 0))
		if err != nil {
			return errors.WithStack(err)
		}

		like, err := ctr.LikeService.Create(user, post)
		if err != nil {
			return errors.WithStack(err)
		}

		return ctx.JSON(http.StatusCreated, map[string]*models.Like{
			"like": like,
		})
	}
}
func LikesDelete(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*models.User)

		like, err := ctr.LikeService.FindById(utils.StringToID(ctx.Param("likeId"), 0))
		if err != nil {
			return errors.WithStack(err)
		}

		if like.UserID != user.ID {
			return ctx.JSON(http.StatusForbidden, map[string]string{
				"message": "You do not have permission to perform this action.",
			})
		}

		if err = ctr.LikeService.Delete(like.ID); err != nil {
			return errors.WithStack(err)
		}

		return ctx.NoContent(http.StatusNoContent)
	}
}
