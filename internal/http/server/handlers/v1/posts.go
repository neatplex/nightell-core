package v1

import (
	"github.com/cockroachdb/errors"
	"github.com/labstack/echo/v4"
	"github.com/neatplex/nightell-core/internal/container"
	"github.com/neatplex/nightell-core/internal/models"
	"github.com/neatplex/nightell-core/internal/utils"
	"net/http"
)

func PostsIndex(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user, err := ctr.UserService.FindBy("id", utils.StringToID(ctx.Param("userId"), 0))
		if err != nil {
			return errors.WithStack(err)
		}
		if user == nil {
			return ctx.NoContent(http.StatusNotFound)
		}

		posts, err := ctr.PostService.Index(
			user.Id,
			utils.StringToID(ctx.QueryParams().Get("lastId"), ^uint64(0)),
			utils.StringToInt(ctx.QueryParams().Get("count"), 100, 10),
		)
		if err != nil {
			return errors.WithStack(err)
		}

		return ctx.JSON(http.StatusCreated, map[string]interface{}{
			"posts": posts,
		})
	}
}

type postsStoreRequest struct {
	Title       string  `json:"title" validate:"required,min=1,max=30"`
	Description string  `json:"description" validate:"max=255"`
	AudioID     uint64  `json:"audio_id" validate:"required"`
	ImageID     *uint64 `json:"image_id"`
}

func PostsStore(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*models.User)

		var r postsStoreRequest
		if err := ctx.Bind(&r); err != nil {
			return err
		}
		if err := ctx.Validate(r); err != nil {
			return err
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
		if s, _ := ctr.PostService.FindBy("audio_id", audio.Id); s != nil {
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
			if s, _ := ctr.PostService.FindBy("image_id", image.Id); s != nil {
				return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{
					"message": "The selected file is already in use.",
				})
			}
			if s, _ := ctr.UserService.FindBy("image_id", image.Id); s != nil {
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
			imageId = &image.Id
		}

		id, err := ctr.PostService.Create(&models.Post{
			UserId:        user.Id,
			Title:         r.Title,
			Description:   r.Description,
			AudioId:       audio.Id,
			ImageId:       imageId,
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

func PostsShow(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		post, err := ctr.PostService.FindById(utils.StringToID(ctx.Param("postId"), 0))
		if err != nil {
			return errors.WithStack(err)
		}
		if post == nil {
			return ctx.NoContent(http.StatusNotFound)
		}
		return ctx.JSON(http.StatusCreated, map[string]interface{}{
			"post": post,
		})
	}
}

type postsUpdateCaptionRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=30"`
	Description string `json:"description" validate:"max=255"`
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

		if post.UserId != user.Id {
			return ctx.JSON(http.StatusForbidden, map[string]string{
				"message": "You do not have permission to perform this action.",
			})
		}

		var r postsUpdateCaptionRequest
		if err = ctx.Bind(&r); err != nil {
			return err
		}
		if err = ctx.Validate(r); err != nil {
			return err
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

		if post.UserId != user.Id {
			return ctx.JSON(http.StatusForbidden, map[string]string{
				"message": "You do not have permission to perform this action.",
			})
		}

		ctr.PostService.Delete(post)

		return ctx.NoContent(http.StatusNoContent)
	}
}
