package v1

import (
	"github.com/cockroachdb/errors"
	"github.com/labstack/echo/v4"
	"github.com/neatplex/nightell-core/internal/models"
	"github.com/neatplex/nightell-core/internal/services/container"
	"github.com/neatplex/nightell-core/internal/utils"
	"net/http"
	"strconv"
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

		followersCount, err := ctr.FollowshipService.CountFollowers(user.ID)
		if err != nil {
			return errors.WithStack(err)
		}

		followingsCount, err := ctr.FollowshipService.CountFollowings(user.ID)
		if err != nil {
			return errors.WithStack(err)
		}

		return ctx.JSON(http.StatusOK, map[string]interface{}{
			"user":             user,
			"followers_count":  followersCount,
			"followings_count": followingsCount,
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
			user.ID,
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
			user.ID,
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

func UsersFollowingsStore(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*models.User)

		if strconv.FormatUint(user.ID, 10) == ctx.Param("followeeId") {
			return ctx.JSON(http.StatusForbidden, map[string]interface{}{
				"message": "You cannot follow yourself!",
			})
		}

		followee, err := ctr.UserService.FindBy("id", utils.StringToID(ctx.Param("followeeId"), 0))
		if err != nil {
			return errors.WithStack(err)
		}
		if followee == nil {
			return ctx.NoContent(http.StatusNotFound)
		}

		if _, err = ctr.FollowshipService.Create(followee.ID, user.ID); err != nil {
			return errors.WithStack(err)
		}

		return ctx.NoContent(http.StatusCreated)
	}
}

func UsersFollowingsDelete(ctr *container.Container) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*models.User)

		if strconv.FormatUint(user.ID, 10) == ctx.Param("followeeId") {
			return ctx.JSON(http.StatusForbidden, map[string]interface{}{
				"message": "You cannot follow yourself!",
			})
		}

		followship, err := ctr.FollowshipService.FindByIds(
			user.ID, utils.StringToID(ctx.Param("followeeId"), 0),
		)
		if err != nil {
			return errors.WithStack(err)
		}
		if followship == nil {
			return ctx.NoContent(http.StatusNotFound)
		}

		err = ctr.FollowshipService.Delete(followship.ID)
		if err != nil {
			return errors.WithStack(err)
		}

		return ctx.NoContent(http.StatusNoContent)
	}
}
