package delivery

import (
	"database/sql"
	"fmt"
	"net/http"
	"technopark_db_forum/internal/models"
	"technopark_db_forum/internal/users/usecase"
	e "technopark_db_forum/pkg/errors"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userUsecase usecase.UsersUsecase
}

func NewUserHandler(userUsecase usecase.UsersUsecase) UserHandler {
	return UserHandler{
		userUsecase: userUsecase,
	}
}

func (h UserHandler) CreateUser(c echo.Context) error {
	nickname := c.Param("nickname")

	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	user.Nickname = nickname

	users, err := h.userUsecase.CreateUser(user)
	if err != nil {
		if err == e.ErrDuplicate {
			return c.JSON(http.StatusConflict, users)
		}
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, users[0])
}

func (h UserHandler) GetUser(c echo.Context) error {
	nickname := c.Param("nickname")

	user, err := h.userUsecase.GetUserByNickname(nickname)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find user by nickname: %s", nickname))
	}
	return c.JSON(http.StatusOK, user)
}

func (h UserHandler) UpdateUser(c echo.Context) error {
	nickname := c.Param("nickname")

	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	user.Nickname = nickname

	updatedUser, err := h.userUsecase.UpdateUser(user)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find user by nickname: %s", nickname))
		}
		if err == e.ErrConflictEmail {
			return echo.NewHTTPError(http.StatusConflict, fmt.Sprintf("This email is already registered by user: %s", updatedUser.Nickname))
		}
		if err == e.ErrConflictNickname {
			return echo.NewHTTPError(http.StatusConflict, fmt.Sprintf("This nickname is already registered by user: %s", updatedUser.Nickname))
		}
		return err
	}
	return c.JSON(http.StatusOK, updatedUser)
}
