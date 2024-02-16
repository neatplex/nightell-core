package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/neatplex/nightel-core/internal/models"
	"github.com/neatplex/nightel-core/internal/services/container"
	"github.com/neatplex/nightel-core/internal/utils"
	"net/http"
	"strconv"
)

func Feed(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*models.User)

		lastId := ^uint64(0)
		requestLastId := ctx.QueryParams().Get("lastId")
		if requestLastId != "" {
			post, err := ctr.PostService.FindById(utils.StringToID(requestLastId, 0))
			if err != nil {
				return err
			}
			if post != nil {
				lastId = post.ID
			}
		}

		count := 10
		requestCount := ctx.QueryParams().Get("count")
		if requestCount != "" {
			parsedRequestCount, _ := strconv.Atoi(requestCount)
			if parsedRequestCount > 0 && parsedRequestCount < 100 {
				count = parsedRequestCount
			}
		}

		posts, err := ctr.PostService.Feed(user.ID, lastId, count)
		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusOK, struct {
			Posts []*models.Post `json:"posts"`
		}{
			Posts: posts,
		})
	}
}
