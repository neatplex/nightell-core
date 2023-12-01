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

type StoriesStoreRequest struct {
	Caption string `json:"caption" validate:"required"`
	Audio   struct {
		Path string `json:"path"`
	} `json:"audio"`
	Image struct {
		Path string `json:"path"`
	} `json:"image"`
	IsPublished bool `json:"is_published"`
}

func StoriesStore(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*models.User)

		var r StoriesStoreRequest
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

		audio, err := ctr.FileService.FindByPath(r.Audio.Path)
		if err != nil {
			return err
		}
		if audio == nil {
			return ctx.JSON(http.StatusNotFound, map[string]string{
				"message": "Audio file not found.",
			})
		}

		var imageId *uint64
		if r.Image.Path != "" {
			image, err := ctr.FileService.FindByPath(r.Image.Path)
			if err != nil {
				return err
			}
			if image == nil {
				return ctx.JSON(http.StatusNotFound, map[string]string{
					"message": "Image file not found.",
				})
			}
			imageId = &image.ID
		}

		err = ctr.StoryService.Create(&models.Story{
			UserID:      user.ID,
			Caption:     r.Caption,
			AudioID:     audio.ID,
			ImageID:     imageId,
			IsPublished: r.IsPublished,
		})
		if err != nil {
			return err
		}

		return ctx.NoContent(http.StatusCreated)
	}
}
