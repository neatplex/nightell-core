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
		Path string `json:"path" validate:"required"`
	} `json:"audio" validate:"required"`
	Image struct {
		Path string `json:"path"`
	} `json:"image"`
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
		audioType, err := ctr.FileService.ToType(audio.Extension)
		if err != nil {
			return err
		}
		if audioType != models.FileTypeAudio {
			return ctx.JSON(http.StatusNotFound, map[string]string{
				"message": "The selected file is not an audio.",
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
			imageType, err := ctr.FileService.ToType(image.Extension)
			if err != nil {
				return err
			}
			if imageType != models.FileTypeImage {
				return ctx.JSON(http.StatusNotFound, map[string]string{
					"message": "The selected file is not an image.",
				})
			}
			imageId = &image.ID
		}

		identity, err := ctr.StoryService.Create(&models.Story{
			UserID:  user.ID,
			Caption: r.Caption,
			AudioID: audio.ID,
			ImageID: imageId,
		})
		if err != nil {
			return err
		}

		story, err := ctr.StoryService.FindByIdentity(identity)
		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusCreated, story)
	}
}

type StoriesUpdateCaptionRequest struct {
	Caption string `json:"caption" validate:"required"`
}

func StoriesUpdateCaption(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*models.User)

		var r StoriesUpdateCaptionRequest
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

		story, err := ctr.StoryService.FindByIdentity(ctx.Param("identity"))
		if err != nil {
			return err
		}
		if story == nil {
			return ctx.NoContent(http.StatusNotFound)
		}

		if story.UserID != user.ID {
			return ctx.JSON(http.StatusForbidden, map[string]string{
				"message": "You do not have permission to perform this action.",
			})
		}

		s := ctr.StoryService.UpdateCaption(story, r.Caption)

		return ctx.JSON(http.StatusCreated, s)
	}
}

func StoriesDelete(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*models.User)

		var r StoriesUpdateCaptionRequest
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

		story, err := ctr.StoryService.FindByIdentity(ctx.Param("identity"))
		if err != nil {
			return err
		}
		if story == nil {
			return ctx.NoContent(http.StatusNotFound)
		}

		if story.UserID != user.ID {
			return ctx.JSON(http.StatusForbidden, map[string]string{
				"message": "You do not have permission to perform this action.",
			})
		}

		ctr.StoryService.Delete(story)

		return ctx.NoContent(http.StatusNoContent)
	}
}
