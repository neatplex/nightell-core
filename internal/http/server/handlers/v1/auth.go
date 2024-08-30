package v1

import (
	"github.com/cockroachdb/errors"
	"github.com/labstack/echo/v4"
	"github.com/neatplex/nightell-core/internal/container"
	"github.com/neatplex/nightell-core/internal/models"
	"github.com/neatplex/nightell-core/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/idtoken"
	"net/http"
)

type signUpRequest struct {
	Username string `json:"username" validate:"required,min=1,max=20"`
	Email    string `json:"email" validate:"required,email,max=191"`
	Password string `json:"password" validate:"required,min=8"`
}

type signInEmailRequest struct {
	Email    string `json:"email" validate:"required,min=1,max=191"`
	Password string `json:"password" validate:"required,min=1"`
}

type signInUsernameRequest struct {
	Username string `json:"username" validate:"required,min=1,max=191"`
	Password string `json:"password" validate:"required,min=1"`
}

type signInGoogleRequest struct {
	Token string `json:"google_token" validate:"required"`
}

type otpEmailSendRequest struct {
	Email string `json:"email" validate:"required,email,max=191"`
}

type otpEmailVerifyRequest struct {
	Email string `json:"email" validate:"required,email,max=191"`
	Otp   string `json:"otp" validate:"required"`
}

func AuthSignUp(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var r signUpRequest
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

		u, err := ctr.UserService.FindBy("username", r.Username)
		if err != nil {
			return errors.WithStack(err)
		}
		if u != nil {
			return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{
				"message": "This username is already reserved.",
			})
		}

		u, err = ctr.UserService.FindBy("email", r.Email)
		if err != nil {
			return errors.WithStack(err)
		}
		if u != nil {
			return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{
				"message": "This email is already registered.",
			})
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
		if err != nil {
			return errors.WithStack(err)
		}

		err = ctr.UserService.Create(&models.User{
			Username: r.Username,
			Email:    r.Email,
			IsBanned: false,
			Password: string(hashedPassword),
		})
		if err != nil {
			return errors.WithStack(err)
		}

		user, err := ctr.UserService.FindBy("email", r.Email)
		if err != nil {
			return errors.WithStack(err)
		}

		token, err := ctr.TokenService.Create(user)
		if err != nil {
			return errors.WithStack(err)
		}

		return ctx.JSON(http.StatusCreated, map[string]interface{}{
			"user":  user,
			"token": token,
		})
	}
}

func AuthOtpEmailSend(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var r otpEmailSendRequest
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

		ttl := ctr.OtpService.Email(r.Email)

		return ctx.JSON(http.StatusCreated, map[string]interface{}{
			"ttl": ttl,
		})
	}
}

func AuthOtpEmailVerify(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var r otpEmailVerifyRequest
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
				"message": "The OTP (one-time password) is not valid.",
			})
		}

		user, err := ctr.UserService.FindBy("email", r.Email)
		if err != nil {
			return errors.WithStack(err)
		}

		if user != nil {
			token, err := ctr.TokenService.Create(user)
			if err != nil {
				return errors.WithStack(err)
			}

			return ctx.JSON(http.StatusCreated, map[string]interface{}{
				"user":  user,
				"token": token,
			})
		} else {
			err = ctr.UserService.Create(&models.User{
				Username: r.Email,
				Email:    r.Email,
				IsBanned: false,
				Password: "",
			})
			if err != nil {
				return errors.WithStack(err)
			}

			user, err = ctr.UserService.FindBy("email", r.Email)
			if err != nil {
				return errors.WithStack(err)
			}

			token, err := ctr.TokenService.Create(user)
			if err != nil {
				return errors.WithStack(err)
			}

			return ctx.JSON(http.StatusCreated, map[string]interface{}{
				"user":  user,
				"token": token,
			})
		}
	}
}

func AuthSignInEmail(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var r signInEmailRequest
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
			if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(r.Password)); err == nil {
				if user.IsBanned {
					return ctx.JSON(http.StatusForbidden, map[string]interface{}{
						"message": "Your account is banned.",
					})
				}

				token, err := ctr.TokenService.FindOrCreate(user)
				if err != nil {
					return errors.WithStack(err)
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

func AuthSignInUsername(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var r signInUsernameRequest
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

		user, err := ctr.UserService.FindBy("username", r.Username)
		if err != nil {
			return errors.WithStack(err)
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
					return errors.WithStack(err)
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

func AuthSignInGoogle(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var r signInGoogleRequest
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

		payload, err := idtoken.Validate(ctx.Request().Context(), r.Token, ctr.Config.Google.OAuthClientId)
		if err != nil {
			return ctx.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Cannot fetch account from Google.",
				"details": err.Error(),
			})
		}

		if !payload.Claims["email_verified"].(bool) {
			return ctx.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Your email is not verified by Google.",
			})
		}

		email := payload.Claims["email"].(string)
		user, err := ctr.UserService.FindBy("email", email)
		if err != nil {
			return errors.WithStack(err)
		}

		if user != nil {
			token, err := ctr.TokenService.Create(user)
			if err != nil {
				return errors.WithStack(err)
			}

			return ctx.JSON(http.StatusCreated, map[string]interface{}{
				"user":  user,
				"token": token,
			})
		} else {
			err = ctr.UserService.Create(&models.User{
				Username: email,
				Email:    email,
				IsBanned: false,
				Password: "",
			})
			if err != nil {
				return errors.WithStack(err)
			}

			user, err = ctr.UserService.FindBy("email", email)
			if err != nil {
				return errors.WithStack(err)
			}

			token, err := ctr.TokenService.Create(user)
			if err != nil {
				return errors.WithStack(err)
			}

			return ctx.JSON(http.StatusCreated, map[string]interface{}{
				"user":  user,
				"token": token,
			})
		}
	}
}
