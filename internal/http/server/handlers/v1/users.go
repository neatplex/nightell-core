package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/neatplex/nightel-core/internal/models"
	"github.com/neatplex/nightel-core/internal/services/container"
	"github.com/neatplex/nightel-core/internal/utils"
	"net/http"
)

func UsersShow(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user, err := ctr.UserService.FindById(utils.StringToID(ctx.Param("userId"), 0))
		if err != nil {
			return err
		}
		if user == nil {
			return ctx.NoContent(http.StatusNotFound)
		}

		return ctx.JSON(http.StatusOK, map[string]*models.User{
			"user": user,
		})
	}
}
