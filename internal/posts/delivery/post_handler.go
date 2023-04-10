package delivery

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"technopark_db_forum/internal/models"
	"technopark_db_forum/internal/posts/usecase"
	e "technopark_db_forum/pkg/errors"

	"github.com/labstack/echo/v4"
)

type PostHandler struct {
	postUsecase usecase.PostUsecase
}

func NewPostHandler(postUsecase usecase.PostUsecase) PostHandler {
	return PostHandler{
		postUsecase: postUsecase,
	}
}

func (h PostHandler) CreatePosts(c echo.Context) error {
	slugOrID := c.Param("slug_or_id")

	var posts []models.Post
	if err := c.Bind(&posts); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	createdPosts, err := h.postUsecase.CreatePosts(posts, slugOrID)
	if err != nil {
		switch err.Error() {
		case e.ErrConflict.Error():
			return c.JSON(http.StatusConflict, createdPosts)
		case e.ErrNotFound.Error():
			return c.JSON(http.StatusNotFound, err)
		case e.ErrOtherThread.Error():
			return echo.NewHTTPError(http.StatusConflict, fmt.Sprintf("Parent post was created in another thread"))
		case e.ErrNoAuthorPost.Error():
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find post author by nickname: %s", posts[0].Author))
		case e.ErrNoRowsID.Error():
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find post thread by id: %s", slugOrID))
		case e.ErrNoRowsSlug.Error():
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find post thread by slug: %s", slugOrID))
		default:
			return c.JSON(http.StatusInternalServerError, err)
		}
	}

	return c.JSON(http.StatusCreated, createdPosts)
}

func (h PostHandler) GetPost(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	related := strings.Split(c.QueryParam("related"), ",")

	post, err := h.postUsecase.GetPostByIDRelared(id, related)
	if err != nil {
		if err.Error() == e.ErrNotFound.Error() {
			return c.JSON(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find post with id: %d", id))
	}

	return c.JSON(http.StatusOK, post)
}

func (h PostHandler) UpdatePost(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	var post models.Post
	if err := c.Bind(&post); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	post.ID = id

	updatedPost, err := h.postUsecase.UpdatePost(post)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return c.JSON(http.StatusNotFound, e.ErrNotFound)
		}
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, updatedPost)
}

func (h PostHandler) GetThreadPosts(c echo.Context) error {
	slugOrID := c.Param("slug_or_id")
	limit, err := strconv.ParseUint(c.QueryParam("limit"), 10, 64)
	if err != nil {
		limit = 100
	}
	since, err := strconv.ParseUint(c.QueryParam("since"), 10, 64)
	if err != nil {
		since = 0
	}
	desc, err := strconv.ParseBool(c.QueryParam("desc"))
	if err != nil {
		desc = false
	}

	sortBy := c.QueryParam("sort")
	if sortBy != "flat" && sortBy != "tree" && sortBy != "parent_tree" {
		sortBy = "flat"
	}

	posts, err := h.postUsecase.GetThreadPosts(slugOrID, limit, sortBy, since, desc)
	if err != nil {
		if err.Error() == e.ErrNotFound.Error() {
			return c.JSON(http.StatusNotFound, err)
		}
		if err.Error() == e.ErrNoRowsID.Error() {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find post thread by id: %s", slugOrID))
		}
		if err.Error() == e.ErrNoRowsSlug.Error() {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find post thread by slug: %s", slugOrID))
		}
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, posts)
}
