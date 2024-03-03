package v1

import (
	"github.com/cockroachdb/errors"
	"github.com/labstack/echo/v4"
	"github.com/neatplex/nightell-core/internal/models"
	"github.com/neatplex/nightell-core/internal/services/container"
	userService "github.com/neatplex/nightell-core/internal/services/user"
	"github.com/neatplex/nightell-core/internal/utils"
	"net/http"
)

func ProfileShow(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*models.User)

		followersCount, err := ctr.FollowshipService.CountFollowers(user.ID)
		if err != nil {
			return errors.WithStack(err)
		}

		followingsCount, err := ctr.FollowshipService.CountFollowings(user.ID)
		if err != nil {
			return errors.WithStack(err)
		}

		u, err := ctr.UserService.FindBy("id", user.ID)
		if err != nil {
			return errors.WithStack(err)
		}

		return ctx.JSON(http.StatusOK, map[string]interface{}{
			"user":             u,
			"followers_count":  followersCount,
			"followings_count": followingsCount,
		})
	}
}

type profileUpdateNameRequest struct {
	Name string `json:"name"`
}

func ProfileUpdateName(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*models.User)

		var r profileUpdateNameRequest
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

		_ = ctr.UserService.UpdateName(user, r.Name)

		u, err := ctr.UserService.FindBy("id", user.ID)
		if err != nil {
			return errors.WithStack(err)
		}

		return ctx.JSON(http.StatusOK, map[string]models.User{
			"user": *u,
		})
	}
}

type profileUpdateBioRequest struct {
	Bio string `json:"bio"`
}

func ProfileUpdateBio(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*models.User)

		var r profileUpdateBioRequest
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

		_ = ctr.UserService.UpdateBio(user, r.Bio)

		u, err := ctr.UserService.FindBy("id", user.ID)
		if err != nil {
			return errors.WithStack(err)
		}

		return ctx.JSON(http.StatusOK, map[string]models.User{
			"user": *u,
		})
	}
}

type profileUpdateUsernameRequest struct {
	Username string `json:"username"`
}

func ProfileUpdateUsername(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*models.User)

		var r profileUpdateUsernameRequest
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

		if !utils.ValidateUsername(r.Username) {
			return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{
				"message": "The username is not valid.",
			})
		}

		_, err := ctr.UserService.UpdateUsername(user, r.Username)
		if err != nil {
			if errors.Is(err, userService.ErrUsernameAlreadyExist) {
				return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{
					"message": "Username already exist.",
				})
			} else {
				return errors.WithStack(err)
			}
		}

		u, err := ctr.UserService.FindBy("id", user.ID)
		if err != nil {
			return errors.WithStack(err)
		}

		return ctx.JSON(http.StatusOK, map[string]models.User{
			"user": *u,
		})
	}
}

type profileUpdateImageRequest struct {
	ImageID *uint64 `json:"image_id"`
}

func ProfileUpdateImage(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*models.User)

		var r profileUpdateImageRequest
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

		image, err := ctr.FileService.FindByID(*r.ImageID)
		if err != nil {
			return errors.WithStack(err)
		}
		if image == nil {
			return ctx.JSON(http.StatusNotFound, map[string]string{
				"message": "Image file not found.",
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

		if s, _ := ctr.PostService.FindBy("image_id", image.ID); s != nil {
			return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{
				"message": "The selected file is already in use.",
			})
		}
		if s, _ := ctr.UserService.FindBy("image_id", image.ID); s != nil {
			return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{
				"message": "The selected file is already in use.",
			})
		}

		_ = ctr.UserService.UpdateImage(user, image.ID)

		u, err := ctr.UserService.FindBy("id", user.ID)
		if err != nil {
			return errors.WithStack(err)
		}

		return ctx.JSON(http.StatusOK, map[string]*models.User{
			"user": u,
		})
	}
}
