package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/neatplex/nightel-core/internal/models"
	"github.com/neatplex/nightel-core/internal/services/container"
	"net/http"
)

func FilesStore(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*models.User)

		fileType, err := ctr.FileService.ToFileType(ctx.FormValue("type"))
		if err != nil {
			return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{
				"message": "File type is not supported.",
			})
		}

		formFile, err := ctx.FormFile("file")
		if err != nil {
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

		path, err := ctr.FileService.Upload(fileHandler)
		if err != nil {
			return err
		}

		err = ctr.FileService.Create(&models.File{
			UserID: user.ID,
			Type:   fileType,
			Path:   path,
		})
		if err != nil {
			return err
		}

		file, err := ctr.FileService.FindByPath(path)
		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusOK, file)
	}
}
