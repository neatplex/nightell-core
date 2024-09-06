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
	Email    string `json:"email" validate:"required,email,min=1,max=191"`
	Password string `json:"password" validate:"required,min=1"`
}

type signInUsernameRequest struct {
	Username string `json:"username" validate:"required,min=1,max=191"`
	Password string `json:"password" validate:"required,min=1"`
}

type signInGoogleRequest struct {
	Token string `json:"google_token" validate:"required"`
}

type otpEmailRequest struct {
	Email string `json:"email" validate:"required,email,max=191"`
}

type otpEmailVerificationRequest struct {
	Email string `json:"email" validate:"required,email,max=191"`
	Otp   string `json:"otp" validate:"required"`
}

func AuthSignUp(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var r signUpRequest
		if err := ctx.Bind(&r); err != nil {
			return err
		}
		if err := ctx.Validate(r); err != nil {
			return err
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

		return createSignInResponse(ctx, ctr, user)
	}
}

func AuthOtpEmail(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var r otpEmailRequest
		if err := ctx.Bind(&r); err != nil {
			return err
		}
		if err := ctx.Validate(r); err != nil {
			return err
		}

		ttl := ctr.OtpService.Email(r.Email)

		return ctx.JSON(http.StatusCreated, map[string]interface{}{
			"ttl": ttl,
		})
	}
}

func AuthOtpEmailVerification(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var r otpEmailVerificationRequest
		if err := ctx.Bind(&r); err != nil {
			return err
		}
		if err := ctx.Validate(r); err != nil {
			return err
		}

		if !ctr.OtpService.Check(r.Email, r.Otp) {
			return ctx.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": "The OTP (one-time password) is incorrect.",
			})
		}

		user, err := ctr.UserService.FindByEmailOrCreate(r.Email)
		if err != nil {
			return errors.WithStack(err)
		}

		return createSignInResponse(ctx, ctr, user)
	}
}

func AuthSignInEmail(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var r signInEmailRequest
		if err := ctx.Bind(&r); err != nil {
			return err
		}
		if err := ctx.Validate(r); err != nil {
			return err
		}

		user, err := ctr.UserService.FindBy("email", r.Email)
		if err != nil {
			return errors.WithStack(err)
		}

		if !ctr.UserService.CheckPassword(user, r.Password) {
			return ctx.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Email or password is incorrect.",
			})
		}

		if user.IsBanned {
			return ctx.JSON(http.StatusForbidden, map[string]interface{}{
				"message": "Your account is banned.",
			})
		}

		return createSignInResponse(ctx, ctr, user)
	}
}

func AuthSignInUsername(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var r signInUsernameRequest
		if err := ctx.Bind(&r); err != nil {
			return err
		}
		if err := ctx.Validate(r); err != nil {
			return err
		}

		user, err := ctr.UserService.FindBy("username", r.Username)
		if err != nil {
			return errors.WithStack(err)
		}

		if !ctr.UserService.CheckPassword(user, r.Password) {
			return ctx.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Username or password is incorrect.",
			})
		}

		if user.IsBanned {
			return ctx.JSON(http.StatusForbidden, map[string]interface{}{
				"message": "Your account is banned.",
			})
		}

		return createSignInResponse(ctx, ctr, user)
	}
}

func AuthSignInGoogle(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var r signInGoogleRequest
		if err := ctx.Bind(&r); err != nil {
			return err
		}
		if err := ctx.Validate(r); err != nil {
			return err
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
		user, err := ctr.UserService.FindByEmailOrCreate(email)
		if err != nil {
			return errors.WithStack(err)
		}

		return createSignInResponse(ctx, ctr, user)
	}
}

func createSignInResponse(ctx echo.Context, ctr *container.Container, user *models.User) error {
	token, err := ctr.TokenService.Create(user)
	if err != nil {
		return errors.WithStack(err)
	}

	return ctx.JSON(http.StatusCreated, map[string]interface{}{
		"user":  user,
		"token": token,
	})
}
