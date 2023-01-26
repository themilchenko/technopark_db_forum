package delivery

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"technopark_db_forum/internal/forum/usecase"
	"technopark_db_forum/internal/models"
	e "technopark_db_forum/pkg/errors"

	"github.com/labstack/echo/v4"
)

type ForumHandler struct {
	ForumUsecase usecase.ForumUsecase
}

func NewForumHandler(forumUsecase usecase.ForumUsecase) ForumHandler {
	return ForumHandler{
		ForumUsecase: forumUsecase,
	}
}

func (h ForumHandler) CreateForum(c echo.Context) error {
	var forum models.Forum
	if err := c.Bind(&forum); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	createdForum, err := h.ForumUsecase.CreateForum(forum)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find user by nickname: %s", forum.UserNickname))
		case e.ErrDuplicate:
			return c.JSON(http.StatusConflict, createdForum)
		case e.ErrThreadNotFound:
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find thread author by nickname: %s", forum.UserNickname))
		default:
			return c.JSON(http.StatusInternalServerError, err)
		}
	}
	return c.JSON(http.StatusCreated, createdForum)
}

func (h ForumHandler) GetForum(c echo.Context) error {
	slug := c.Param("slug")

	forum, err := h.ForumUsecase.GetForumBySlug(slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find forum with slug: %s", slug))
		}
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, forum)
}

func (h ForumHandler) GetForumUsers(c echo.Context) error {
	slug := c.Param("slug")

	limit, err := strconv.ParseUint(c.QueryParam("limit"), 10, 64)
	if err != nil {
		limit = 100
	}
	// if limit == 10 {
	// 	fmt.Sprint(limit)
	// }
	since := c.QueryParam("since")
	desc, err := strconv.ParseBool(c.QueryParam("desc"))
	if err != nil {
		desc = false
	}

	users, err := h.ForumUsecase.GetForumUsersBySlug(slug, models.ThreadOptions{
		Limit: limit,
		Since: since,
		Desc:  desc,
	})
	// if len(users) == 0 {
	// 	return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find forum by slug: %s", slug))
	// }
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, err)
		}
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, users)
}
