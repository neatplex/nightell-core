package v1

import (
	"github.com/cockroachdb/errors"
	"github.com/labstack/echo/v4"
	"github.com/neatplex/nightel-core/internal/models"
	"github.com/neatplex/nightel-core/internal/services/container"
	"github.com/neatplex/nightel-core/internal/utils"
	"net/http"
)

func PostsIndex(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user, err := ctr.UserService.FindById(utils.StringToID(ctx.Param("userId"), 0))
		if err != nil {
			return errors.WithStack(err)
		}
		if user == nil {
			return ctx.NoContent(http.StatusNotFound)
		}

		posts, err := ctr.PostService.Index(user.ID)
		if err != nil {
			return errors.WithStack(err)
		}
		return ctx.JSON(http.StatusCreated, map[string]interface{}{
			"posts": posts,
		})
	}
}

type PostsStoreRequest struct {
	Title       string  `json:"title" validate:"required,min=1,max=50"`
	Description string  `json:"description" validate:"max=300"`
	AudioID     uint64  `json:"audio_id" validate:"required"`
	ImageID     *uint64 `json:"image_id"`
}

func PostsStore(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*models.User)

		var r PostsStoreRequest
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

		audio, err := ctr.FileService.FindByID(r.AudioID)
		if err != nil {
			return errors.WithStack(err)
		}
		if audio == nil {
			return ctx.JSON(http.StatusNotFound, map[string]string{
				"message": "Audio file not found.",
			})
		}
		if s, _ := ctr.PostService.FindBy("audio_id", audio.ID); s != nil {
			return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{
				"message": "The selected file is already in use.",
			})
		}
		audioType, err := ctr.FileService.TypeFromExtension(audio.Extension)
		if err != nil {
			return errors.WithStack(err)
		}
		if audioType != models.FileTypeAudio {
			return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{
				"message": "The selected file is not an audio.",
			})
		}

		var imageId *uint64
		if r.ImageID != nil {
			image, err := ctr.FileService.FindByID(*r.ImageID)
			if err != nil {
				return errors.WithStack(err)
			}
			if image == nil {
				return ctx.JSON(http.StatusNotFound, map[string]string{
					"message": "Image file not found.",
				})
			}
			if s, _ := ctr.PostService.FindBy("image_id", image.ID); s != nil {
				return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{
					"message": "The selected file is already in use.",
				})
			}
			imageType, err := ctr.FileService.TypeFromExtension(image.Extension)
			if err != nil {
				return errors.WithStack(err)
			}
			if imageType != models.FileTypeImage {
				return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{
					"message": "The selected file is not an image.",
				})
			}
			imageId = &image.ID
		}

		id, err := ctr.PostService.Create(&models.Post{
			UserID:        user.ID,
			Title:         r.Title,
			Description:   r.Description,
			AudioID:       audio.ID,
			ImageID:       imageId,
			LikesCount:    0,
			CommentsCount: 0,
		})
		if err != nil {
			return errors.WithStack(err)
		}

		post, err := ctr.PostService.FindById(id)
		if err != nil {
			return errors.WithStack(err)
		}

		return ctx.JSON(http.StatusCreated, map[string]interface{}{
			"post": post,
		})
	}
}

type PostsUpdateCaptionRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=50"`
	Description string `json:"description" validate:"max=300"`
}

func PostsUpdate(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*models.User)

		post, err := ctr.PostService.FindById(utils.StringToID(ctx.Param("postId"), 0))
		if err != nil {
			return errors.WithStack(err)
		}
		if post == nil {
			return ctx.NoContent(http.StatusNotFound)
		}

		if post.UserID != user.ID {
			return ctx.JSON(http.StatusForbidden, map[string]string{
				"message": "You do not have permission to perform this action.",
			})
		}

		var r PostsUpdateCaptionRequest
		if err = ctx.Bind(&r); err != nil {
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"message": "Cannot parse the request body.",
			})
		}
		if err = ctx.Validate(r); err != nil {
			return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{
				"message": err.Error(),
			})
		}

		s := ctr.PostService.UpdateFields(post, r.Title, r.Description)

		return ctx.JSON(http.StatusCreated, map[string]interface{}{
			"post": s,
		})
	}
}

func PostsDelete(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*models.User)

		post, err := ctr.PostService.FindById(utils.StringToID(ctx.Param("postId"), 0))
		if err != nil {
			return errors.WithStack(err)
		}
		if post == nil {
			return ctx.NoContent(http.StatusNotFound)
		}

		if post.UserID != user.ID {
			return ctx.JSON(http.StatusForbidden, map[string]string{
				"message": "You do not have permission to perform this action.",
			})
		}

		ctr.PostService.Delete(post)

		return ctx.NoContent(http.StatusNoContent)
	}
}
