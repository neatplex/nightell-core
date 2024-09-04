package v1

import (
	"github.com/cockroachdb/errors"
	"github.com/labstack/echo/v4"
	"github.com/neatplex/nightell-core/internal/container"
	"github.com/neatplex/nightell-core/internal/models"
	userService "github.com/neatplex/nightell-core/internal/services/user"
	"github.com/neatplex/nightell-core/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func ProfileShow(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*models.User)

		followersCount, err := ctr.FollowshipService.CountFollowers(user.Id)
		if err != nil {
			return errors.WithStack(err)
		}

		followingsCount, err := ctr.FollowshipService.CountFollowings(user.Id)
		if err != nil {
			return errors.WithStack(err)
		}

		u, err := ctr.UserService.FindBy("id", user.Id)
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
	Name string `json:"name" validate:"max=20"`
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

		if _, err := ctr.UserService.UpdateName(user, r.Name); err != nil {
			return errors.WithStack(err)
		}

		u, err := ctr.UserService.FindBy("id", user.Id)
		if err != nil {
			return errors.WithStack(err)
		}

		return ctx.JSON(http.StatusOK, map[string]models.User{
			"user": *u,
		})
	}
}

type profileUpdateBioRequest struct {
	Bio string `json:"bio" validate:"max=255"`
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

		if _, err := ctr.UserService.UpdateBio(user, r.Bio); err != nil {
			return errors.WithStack(err)
		}

		u, err := ctr.UserService.FindBy("id", user.Id)
		if err != nil {
			return errors.WithStack(err)
		}

		return ctx.JSON(http.StatusOK, map[string]models.User{
			"user": *u,
		})
	}
}

type profileUpdatePasswordRequest struct {
	Password string `json:"password" validate:"min=8,max=255"`
}

func ProfileUpdatePassword(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*models.User)

		var r profileUpdatePasswordRequest
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

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
		if err != nil {
			return errors.WithStack(err)
		}

		if _, err = ctr.UserService.UpdatePassword(user, string(hashedPassword)); err != nil {
			return errors.WithStack(err)
		}

		u, err := ctr.UserService.FindBy("id", user.Id)
		if err != nil {
			return errors.WithStack(err)
		}

		return ctx.JSON(http.StatusOK, map[string]models.User{
			"user": *u,
		})
	}
}

type profileUpdateUsernameRequest struct {
	Username string `json:"username" validate:"required,min=5,max=20"`
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

		u, err := ctr.UserService.FindBy("id", user.Id)
		if err != nil {
			return errors.WithStack(err)
		}

		return ctx.JSON(http.StatusOK, map[string]models.User{
			"user": *u,
		})
	}
}

type profileUpdateEmailRequest struct {
	Email string `json:"email" validate:"required,email,max=191"`
}

func ProfileUpdateEmail(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*models.User)

		var r profileUpdateEmailRequest
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

		if user.Email == r.Email {
			return ctx.JSON(http.StatusOK, map[string]models.User{
				"user": *user,
			})
		}

		ttl := ctr.OtpService.Email(r.Email)

		return ctx.JSON(http.StatusOK, map[string]interface{}{
			"ttl": ttl,
		})
	}
}

type profileUpdateEmailVerifyRequest struct {
	Email string `json:"email" validate:"required,email,max=191"`
	Otp   string `json:"otp" validate:"required"`
}

func ProfileUpdateEmailVerify(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*models.User)

		var r profileUpdateEmailVerifyRequest
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

		isValid := ctr.OtpService.Check(r.Email, r.Otp)
		if !isValid {
			return ctx.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": "The OTP (one-time password) is incorrect.",
			})
		}

		if user.Email == r.Email {
			return ctx.JSON(http.StatusOK, map[string]models.User{
				"user": *user,
			})
		}

		_, err := ctr.UserService.UpdateEmail(user, r.Email)
		if err != nil {
			if errors.Is(err, userService.ErrEmailAlreadyExist) {
				return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{
					"message": "Email already exist.",
				})
			} else {
				return errors.WithStack(err)
			}
		}

		u, err := ctr.UserService.FindBy("id", user.Id)
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

		if _, err = ctr.UserService.UpdateImage(user, image.Id); err != nil {
			return errors.WithStack(err)
		}

		u, err := ctr.UserService.FindBy("id", user.Id)
		if err != nil {
			return errors.WithStack(err)
		}

		return ctx.JSON(http.StatusOK, map[string]*models.User{
			"user": u,
		})
	}
}

func ProfileDelete(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*models.User)

		if err := ctr.UserService.Delete(user); err != nil {
			return errors.WithStack(err)
		}

		return ctx.NoContent(http.StatusNoContent)
	}
}
