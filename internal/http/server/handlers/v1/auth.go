package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/neatplex/nightel-core/internal/models"
	"github.com/neatplex/nightel-core/internal/services/container"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type SignUpRequest struct {
	Name     string `json:"name" validate:"required,min=1,max=128"`
	Email    string `json:"email" validate:"required,email,max=128"`
	Password string `json:"password" validate:"required,min=8,max=128"`
}

type SignInRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func AuthSignIn(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var r SignInRequest
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

		user, err := ctr.UserService.FindByEmail(r.Email)
		if err != nil {
			return err
		}
		if user != nil {
			if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(r.Password)); err == nil {
				if user.IsBanned {
					return ctx.JSON(http.StatusForbidden, map[string]interface{}{
						"message": "Your account is banned.",
					})
				}

				token, err := ctr.TokenService.FindOrCreate(user)
				if err != nil {
					return err
				}

				return ctx.JSON(http.StatusCreated, map[string]interface{}{
					"user":  user,
					"token": token,
				})
			}
		}

		return ctx.JSON(http.StatusUnauthorized, map[string]string{
			"message": "Email or password is incorrect.",
		})
	}
}

func AuthSignUp(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var r SignUpRequest
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

		u, err := ctr.UserService.FindByEmail(r.Email)
		if err != nil {
			return err
		}
		if u != nil {
			return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{
				"message": "This email is already registered.",
			})
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		err = ctr.UserService.Create(&models.User{
			Name:     r.Name,
			Email:    r.Email,
			IsTeller: false,
			IsBanned: false,
			Password: string(hashedPassword),
		})
		if err != nil {
			return err
		}

		user, err := ctr.UserService.FindByEmail(r.Email)
		if err != nil {
			return err
		}

		token, err := ctr.TokenService.Create(user)
		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusCreated, map[string]interface{}{
			"user":  user,
			"token": token,
		})
	}
}
