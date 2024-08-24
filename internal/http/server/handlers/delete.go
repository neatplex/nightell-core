package handlers

import (
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/labstack/echo/v4"
	"github.com/neatplex/nightell-core/internal/container"
	"net/http"
)

type deleteRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func DeleteRequest(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var r deleteRequest
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

		user, err := ctr.UserService.FindBy("email", r.Email)
		if err != nil {
			return errors.WithStack(err)
		}

		if user != nil {
			code, err := ctr.RemoveService.FindOrCreate(user)
			if err != nil {
				return errors.WithStack(err)
			}
			link := fmt.Sprintf("%s/delete-account?q=%s", ctr.Config.URL, code)
			ctr.Mailer.SendDeleteAccount(user.Email, user.Username, link)
		}

		return ctx.JSON(200, map[string]string{
			"message": "Request submitted successfully.",
		})
	}
}

func DeleteAccount(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		code := ctx.QueryParam("q")
		remove, err := ctr.RemoveService.FindBy("code", code)
		if err != nil {
			return errors.WithStack(err)
		}

		if remove != nil {
			if err = ctr.UserService.Delete(remove.User); err != nil {
				return errors.WithStack(err)
			}
		}

		return ctx.Redirect(http.StatusTemporaryRedirect, "/bye.html")
	}
}
