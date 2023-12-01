package handlers

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Debug(ctx echo.Context) error {
	jsonBody := make(map[string]interface{})
	_ = json.NewDecoder(ctx.Request().Body).Decode(&jsonBody)

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"header": ctx.Request().Header,
		"body":   jsonBody,
		"ip":     ctx.RealIP(),
	})
}
