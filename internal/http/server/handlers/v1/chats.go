package v1

import (
	"github.com/cockroachdb/errors"
	"github.com/labstack/echo/v4"
	"github.com/neatplex/nightell-core/internal/models"
	"github.com/neatplex/nightell-core/internal/services/container"
	"github.com/neatplex/nightell-core/internal/utils"
	"net/http"
)

func ChatsIndex(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*models.User)

		chats, err := ctr.ChatService.ByUser(
			user.ID,
			utils.StringToID(ctx.QueryParams().Get("lastId"), ^uint64(0)),
			utils.StringToInt(ctx.QueryParams().Get("count"), 20),
		)
		if err != nil {
			return errors.WithStack(err)
		}

		return ctx.JSON(http.StatusCreated, map[string]interface{}{
			"chats": chats,
		})
	}
}

type chatStoreRequest struct {
	ToID uint64 `json:"to_id" validate:"required,number"`
}

func ChatsStore(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var r chatStoreRequest
		if err := ctx.Bind(&r); err != nil {
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"message": "Cannot parse the request body.",
			})
		}
		if err := ctx.Validate(r); err != nil {
			return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{
				"message": err.Error(),
			})
		}

		from := ctx.Get("user").(*models.User)

		to, err := ctr.UserService.FindBy("id", r.ToID)
		if err != nil {
			return errors.WithStack(err)
		}
		if to == nil {
			return ctx.NoContent(http.StatusNotFound)
		}

		chat, err := ctr.ChatService.OneByUsers([]uint64{from.ID, to.ID})
		if err != nil {
			return errors.WithStack(err)
		}

		if chat == nil {
			chat = &models.Chat{FromID: from.ID, ToID: to.ID}
			if err = ctr.ChatService.Create(chat); err != nil {
				return errors.WithStack(err)
			}
			if chat, err = ctr.ChatService.OneBy("id", chat.ID); err != nil {
				return errors.WithStack(err)
			}
		}

		return ctx.JSON(http.StatusCreated, map[string]interface{}{
			"chat": chat,
		})
	}
}
