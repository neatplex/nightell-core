package v1

import (
	"github.com/cockroachdb/errors"
	"github.com/labstack/echo/v4"
	"github.com/neatplex/nightell-core/internal/container"
	"github.com/neatplex/nightell-core/internal/models"
	"github.com/neatplex/nightell-core/internal/utils"
	"net/http"
)

func CommentsIndexByUser(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user, err := ctr.UserService.FindBy("id", utils.StringToID(ctx.Param("userId"), 0))
		if err != nil {
			return errors.WithStack(err)
		}
		if user == nil {
			return ctx.NoContent(http.StatusNotFound)
		}
		if user.Id != ctx.Get("user").(*models.User).Id {
			return ctx.NoContent(http.StatusForbidden)
		}

		comments, err := ctr.CommentService.IndexByUser(
			user.Id,
			utils.StringToID(ctx.QueryParams().Get("lastId"), ^uint64(0)),
			utils.StringToInt(ctx.QueryParams().Get("count"), 100, 10),
		)
		if err != nil {
			return errors.WithStack(err)
		}

		return ctx.JSON(http.StatusCreated, map[string]interface{}{
			"comments": comments,
		})
	}
}

func CommentsIndexByPost(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		post, err := ctr.PostService.FindBy("id", utils.StringToID(ctx.Param("postId"), 0))
		if err != nil {
			return errors.WithStack(err)
		}
		if post == nil {
			return ctx.NoContent(http.StatusNotFound)
		}

		comments, err := ctr.CommentService.IndexByPost(
			post.Id,
			utils.StringToID(ctx.QueryParams().Get("lastId"), ^uint64(0)),
			utils.StringToInt(ctx.QueryParams().Get("count"), 100, 10),
		)
		if err != nil {
			return errors.WithStack(err)
		}

		return ctx.JSON(http.StatusCreated, map[string]interface{}{
			"comments": comments,
		})
	}
}

type commentsStoreRequest struct {
	Text   string `json:"text" validate:"min=1,max=255"`
	PostId uint64 `json:"post_id" validate:"required"`
}

func CommentsStore(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*models.User)

		var r commentsStoreRequest
		if err := ctx.Bind(&r); err != nil {
			return err
		}
		if err := ctx.Validate(r); err != nil {
			return err
		}

		post, err := ctr.PostService.FindById(r.PostId)
		if err != nil {
			return errors.WithStack(err)
		}
		if post == nil {
			return ctx.JSON(http.StatusNotFound, map[string]string{
				"message": "Post not found.",
			})
		}

		id, err := ctr.CommentService.Create(&models.Comment{
			PostId: post.Id,
			UserId: user.Id,
			Text:   r.Text,
		})
		if err != nil {
			return errors.WithStack(err)
		}

		comment, err := ctr.CommentService.FindById(id)
		if err != nil {
			return errors.WithStack(err)
		}

		return ctx.JSON(http.StatusCreated, map[string]interface{}{
			"comment": comment,
		})
	}
}

func CommentsDelete(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*models.User)

		comment, err := ctr.CommentService.FindById(utils.StringToID(ctx.Param("commentId"), 0))
		if err != nil {
			return errors.WithStack(err)
		}
		if comment == nil {
			return ctx.NoContent(http.StatusNotFound)
		}

		if comment.UserId != user.Id && comment.Post.UserId != user.Id {
			return ctx.JSON(http.StatusForbidden, map[string]string{
				"message": "You do not have permission to perform this action.",
			})
		}

		ctr.CommentService.Delete(comment)

		return ctx.NoContent(http.StatusNoContent)
	}
}
