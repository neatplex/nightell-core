package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/neatplex/nightel-core/internal/models"
	"github.com/neatplex/nightel-core/internal/services/container"
	"net/http"
)

type LikesStoreRequest struct {
	StoryID uint64 `json:"story_id" validate:"required"`
}

func LikesStore(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*models.User)

		var r LikesStoreRequest
		if err := ctx.Bind(&r); err != nil {
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"message": "Cannot parse the request body.",
			})
		}
		if err := ctx.Validate(r); err != nil {
			return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{
				"message": err.Error(),
			})
		}

		story, err := ctr.StoryService.FindById(r.StoryID)
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

type LikesDeleteRequest struct {
	LikeID uint64 `json:"like_id" validate:"required"`
}

func LikesDelete(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*models.User)

		var r LikesDeleteRequest
		if err := ctx.Bind(&r); err != nil {
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"message": "Cannot parse the request body.",
			})
		}
		if err := ctx.Validate(r); err != nil {
			return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{
				"message": err.Error(),
			})
		}

		like, err := ctr.LikeService.FindById(r.LikeID)
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
