package delivery

import (
	"net/http"
	"technopark_db_forum/internal/service/usecase"

	"github.com/labstack/echo/v4"
)

type ServiceHandler interface {
	GetStatus(c echo.Context) error
	Clear(c echo.Context) error
}

type serviceHandler struct {
	serviceUsecase usecase.ServiceUsecase
}

func NewServiceHandler(serviceUsecase usecase.ServiceUsecase) ServiceHandler {
	return &serviceHandler{
		serviceUsecase: serviceUsecase,
	}
}

func (h serviceHandler) GetStatus(c echo.Context) error {
	status, err := h.serviceUsecase.GetStatus()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, status)
}

func (h serviceHandler) Clear(c echo.Context) error {
	err := h.serviceUsecase.Clear()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, struct{}{})
}
