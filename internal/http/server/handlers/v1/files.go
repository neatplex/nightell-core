package v1

import (
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/labstack/echo/v4"
	"github.com/neatplex/nightell-core/internal/container"
	"github.com/neatplex/nightell-core/internal/models"
	"net/http"
)

func FilesStore(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*models.User)

		extension, err := ctr.FileService.ExtensionFromString(ctx.FormValue("extension"))
		if err != nil {
			return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{
				"message": fmt.Sprintf("Extension ``%s is not supported.", ctx.FormValue("extension")),
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

		path := ctr.FileService.Path(extension)

		err = ctr.FileService.Create(&models.File{
			UserID:    user.ID,
			Extension: extension,
			Path:      path,
		})
		if err != nil {
			return errors.WithStack(err)
		}

		file, err := ctr.FileService.FindByPath(path)
		if err != nil {
			return errors.WithStack(err)
		}

		if err = ctr.FileService.Upload(fileHandler, path); err != nil {
			return errors.WithStack(err)
		}

		return ctx.JSON(http.StatusOK, map[string]models.File{
			"file": *file,
		})
	}
}
