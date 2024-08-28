package v1

import (
	"github.com/cockroachdb/errors"
	"github.com/labstack/echo/v4"
	"github.com/neatplex/nightell-core/internal/container"
	"github.com/neatplex/nightell-core/internal/models"
	"github.com/neatplex/nightell-core/internal/utils"
	"net/http"
)

func UsersShow(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user, err := ctr.UserService.FindBy("id", utils.StringToID(ctx.Param("userId"), 0))
		if err != nil {
			return errors.WithStack(err)
		}
		if user == nil {
			return ctx.NoContent(http.StatusNotFound)
		}

		followersCount, err := ctr.FollowshipService.CountFollowers(user.Id)
		if err != nil {
			return errors.WithStack(err)
		}

		followingsCount, err := ctr.FollowshipService.CountFollowings(user.Id)
		if err != nil {
			return errors.WithStack(err)
		}

		authUser := ctx.Get("user").(*models.User)

		relation, err := ctr.FollowshipService.FindByIds(user.Id, authUser.Id)
		if err != nil {
			return errors.WithStack(err)
		}
		followsMe := relation != nil

		relation, err = ctr.FollowshipService.FindByIds(authUser.Id, user.Id)
		if err != nil {
			return errors.WithStack(err)
		}
		followedByMe := relation != nil

		return ctx.JSON(http.StatusOK, map[string]interface{}{
			"user":             user,
			"followers_count":  followersCount,
			"followings_count": followingsCount,
			"follows_me":       followsMe,
			"followed_by_me":   followedByMe,
		})
	}
}

func UsersFollowers(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user, err := ctr.UserService.FindBy("id", utils.StringToID(ctx.Param("userId"), 0))
		if err != nil {
			return errors.WithStack(err)
		}
		if user == nil {
			return ctx.NoContent(http.StatusNotFound)
		}

		followships, err := ctr.FollowshipService.IndexFollowers(
			user.Id,
			utils.StringToID(ctx.QueryParams().Get("lastId"), ^uint64(0)),
			utils.StringToInt(ctx.QueryParams().Get("count"), 100, 10),
		)
		if err != nil {
			return errors.WithStack(err)
		}

		users := make([]*models.User, 0, len(followships))
		for _, f := range followships {
			users = append(users, f.Follower)
		}

		return ctx.JSON(http.StatusOK, map[string]interface{}{
			"users": users,
		})
	}
}

func UsersFollowings(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user, err := ctr.UserService.FindBy("id", utils.StringToID(ctx.Param("userId"), 0))
		if err != nil {
			return errors.WithStack(err)
		}
		if user == nil {
			return ctx.NoContent(http.StatusNotFound)
		}

		followships, err := ctr.FollowshipService.IndexFollowings(
			user.Id,
			utils.StringToID(ctx.QueryParams().Get("lastId"), ^uint64(0)),
			utils.StringToInt(ctx.QueryParams().Get("count"), 100, 10),
		)
		if err != nil {
			return errors.WithStack(err)
		}

		users := make([]*models.User, 0, len(followships))
		for _, f := range followships {
			users = append(users, f.Followee)
		}

		return ctx.JSON(http.StatusOK, map[string]interface{}{
			"users": users,
		})
	}
}

func UsersFollowersStore(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userId := utils.StringToID(ctx.Param("userId"), 0)
		authUser := ctx.Get("user").(*models.User)

		if authUser.Id == userId {
			return ctx.JSON(http.StatusForbidden, map[string]interface{}{
				"message": "You cannot follow yourself!",
			})
		}

		followee, err := ctr.UserService.FindBy("id", userId)
		if err != nil {
			return errors.WithStack(err)
		}
		if followee == nil {
			return ctx.NoContent(http.StatusNotFound)
		}

		if _, err = ctr.FollowshipService.Create(followee.Id, authUser.Id); err != nil {
			return errors.WithStack(err)
		}

		return ctx.NoContent(http.StatusCreated)
	}
}

func UsersFollowersDelete(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userId := utils.StringToID(ctx.Param("userId"), 0)
		authUser := ctx.Get("user").(*models.User)

		followship, err := ctr.FollowshipService.FindByIds(authUser.Id, userId)
		if err != nil {
			return errors.WithStack(err)
		}
		if followship == nil {
			return ctx.NoContent(http.StatusNoContent)
		}

		err = ctr.FollowshipService.Delete(followship.Id)
		if err != nil {
			return errors.WithStack(err)
		}

		return ctx.NoContent(http.StatusNoContent)
	}
}
