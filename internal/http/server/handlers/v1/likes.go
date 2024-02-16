package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/neatplex/nightel-core/internal/models"
	"github.com/neatplex/nightel-core/internal/services/container"
	"github.com/neatplex/nightel-core/internal/utils"
	"net/http"
)

func LikesIndex(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		stories, err := ctr.LikeService.IndexByStoryIDWithUser(
			utils.StringToID(ctx.Param("storyId"), 0),
			utils.StringToID(ctx.QueryParams().Get("lastId"), ^uint64(0)),
			utils.StringToInt(ctx.QueryParams().Get("count"), 10),
		)
		if err != nil {
			return err
		}
		return ctx.JSON(http.StatusCreated, map[string]interface{}{
			"likes": stories,
		})
	}
}

func LikesStore(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*models.User)

		story, err := ctr.StoryService.FindById(utils.StringToID(ctx.Param("storyId"), 0))
		if err != nil {
			return err
		}

		like, err := ctr.LikeService.Create(user, story)
		if err != nil {
			return err
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
			return err
		}

		if like.UserID != user.ID {
			return ctx.JSON(http.StatusForbidden, map[string]string{
				"message": "You do not have permission to perform this action.",
			})
		}

		if err = ctr.LikeService.Delete(like.ID); err != nil {
			return err
		}

		return ctx.NoContent(http.StatusNoContent)
	}
}
