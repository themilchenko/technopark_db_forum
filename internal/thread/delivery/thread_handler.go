package delivery

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"technopark_db_forum/internal/models"
	"technopark_db_forum/internal/thread/usecase"
	e "technopark_db_forum/pkg/errors"
	"time"

	"github.com/labstack/echo/v4"
)

type ThreadHandler struct {
	threadUsecase usecase.ThreadUsecase
}

func NewThreadHandler(threadUsecase usecase.ThreadUsecase) ThreadHandler {
	return ThreadHandler{
		threadUsecase: threadUsecase,
	}
}

// thread/slug/create
func (h ThreadHandler) CreateThread(c echo.Context) error {
	slug := c.Param("slug")

	var thread models.Thread
	if err := c.Bind(&thread); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	thread.Forum = slug

	createdThread, err := h.threadUsecase.CreateThread(thread)
	if err != nil {
		if err == e.ErrDuplicate {
			return c.JSON(http.StatusConflict, createdThread)
		} else if err == e.ErrConflictNickname {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find thread author by nickname: %s", thread.Author))
		} else if err == e.ErrThreadNotFound {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find thread forum by slug: %s", thread.Forum))
		}
		return err
	}
	return c.JSON(http.StatusCreated, createdThread)
}

// forum/slug/threads
func (h ThreadHandler) GetThreadMsgs(c echo.Context) error {
	slugOrID := c.Param("slug")
	var threadOptions models.ThreadOptions
	limit, err := strconv.ParseUint(c.QueryParam("limit"), 10, 64)
	if err != nil {
		threadOptions.Limit = 100
	}
	threadOptions.Limit = limit

	var since time.Time
	if c.QueryParam("since") != "" {
		s := c.QueryParam("since")
		since, err = time.Parse(time.RFC3339, s)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
	}
	desk, err := strconv.ParseBool(c.QueryParam("desc"))
	if err != nil {
		threadOptions.Desc = false
	}
	threadOptions.Desc = desk

	thread, err := h.threadUsecase.GetThreadMsgsBySlug(slugOrID, since, threadOptions)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusOK, thread)
}

func (h ThreadHandler) UpdateThread(c echo.Context) error {
	slugOrID := c.Param("slug_or_id")

	var thread models.Thread
	if err := c.Bind(&thread); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	// thread.Slug = slugOrID

	updatedThread, err := h.threadUsecase.UpdateThread(thread, slugOrID)
	if err != nil {
		if errors.Is(err, e.ErrNoRowsID) {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find thread with id: %d", thread.ID))
		}
		if errors.Is(err, e.ErrNoRowsSlug) {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find thread with slug: %s", thread.Slug))
		}
		return err
	}
	return c.JSON(http.StatusOK, updatedThread)
}

func (h ThreadHandler) GetThread(c echo.Context) error {
	slugOrID := c.Param("slug_or_id")

	thread, err := h.threadUsecase.GetThread(slugOrID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, err)
		}
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, thread)
}

func (h ThreadHandler) CreateVote(c echo.Context) error {
	slugOrID := c.Param("slug_or_id")

	vote := models.Vote{}

	if err := c.Bind(&vote); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	response, err := h.threadUsecase.CreateVote(vote, slugOrID)
	if err == sql.ErrNoRows {
		if response.ID == 0 {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find thread with slug: %s", slugOrID))
		} else if response.Slug == "" {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find thread with id: %d", response.ID))
		}

		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find thread with slug: %s", slugOrID))
	} else if err == e.ErrUserNotFound {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find user by nickname: %s", vote.Nickname))
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, response)


	// slugOrID := c.Param("slug_or_id")

	// var vote models.Vote
	// if err := c.Bind(&vote); err != nil {
	// 	return c.JSON(http.StatusBadRequest, err)
	// }
	// if vote.VoiceValue == -1 {
	// 	fmt.Println("voice value -1")
	// }

	// updatedThread, err := h.threadUsecase.CreateVote(vote, slugOrID)
	// if err != nil {
	// 	if errors.Is(err, sql.ErrNoRows) {
	// 		return c.JSON(http.StatusNotFound, err)
	// 	}
	// 	return err
	// }
	// return c.JSON(http.StatusOK, updatedThread)
}
