package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Index(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "It is the homepage!",
	})
}
