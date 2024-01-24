package v1

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/neatplex/nightel-core/internal/logger"
	"github.com/neatplex/nightel-core/internal/models"
	"github.com/neatplex/nightel-core/internal/services/container"
	"go.uber.org/zap"
	"net/http"
)

func FilesStore(ctr *container.Container, l *logger.Logger) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*models.User)

		extension, err := ctr.FileService.ToExtension(ctx.FormValue("extension"))
		if err != nil {
			return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{
				"message": fmt.Sprintf("Extension ``%s is not supported.", ctx.FormValue("extension")),
			})
		}

		formFile, err := ctx.FormFile("file")
		if err != nil {
			l.Debug("cannot get form file", zap.Error(err))
			return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{
				"message": "File is not uploaded.",
			})
		}

		fileHandler, err := formFile.Open()
		if err != nil {
			return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{
				"message": "File is not processable.",
			})
		}

		path, err := ctr.FileService.Upload(fileHandler, extension)
		if err != nil {
			return err
		}

		err = ctr.FileService.Create(&models.File{
			UserID:    user.ID,
			Extension: extension,
			Path:      path,
		})
		if err != nil {
			return err
		}

		file, err := ctr.FileService.FindByPath(path)
		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusOK, map[string]models.File{
			"file": *file,
		})
	}
}
