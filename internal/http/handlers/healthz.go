package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func Healthz(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNoContent)
}
