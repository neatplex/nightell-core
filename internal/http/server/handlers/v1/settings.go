package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/neatplex/nightell-core/internal/container"
	"net/http"
)

func SettingsIndex(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, *ctr.SettingService.Get())
	}
}
