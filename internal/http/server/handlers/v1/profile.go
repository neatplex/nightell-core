package v1

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/neatplex/nightel-core/internal/models"
	"github.com/neatplex/nightel-core/internal/services/container"
	userService "github.com/neatplex/nightel-core/internal/services/user"
	"net/http"
)

func ProfileShow() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*models.User)

		return ctx.JSON(http.StatusOK, user)
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

		u := ctr.UserService.UpdateName(user, r.Name)

		return ctx.JSON(http.StatusOK, u)
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

		u := ctr.UserService.UpdateBio(user, r.Bio)

		return ctx.JSON(http.StatusOK, u)
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

		u, err := ctr.UserService.UpdateUsername(user, r.Username)
		if err != nil {
			if errors.Is(err, userService.ErrUsernameAlreadyExist) {
				return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{
					"message": "Username already exist.",
				})
			} else {
				return err
			}
		}

		return ctx.JSON(http.StatusOK, u)
	}
}
