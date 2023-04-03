package main

import (
	"technopark_db_forum/internal/app"

	"github.com/labstack/echo/v4"
)

const (
	defaultDBConfig = "host=localhost port=5432 dbname=dev sslmode=disable"
)

func main() {
	e := echo.New()
	s := app.New(e)

	if err := s.Start("0.0.0.0:8080", defaultDBConfig); err != nil {
		s.Echo.Logger.Error("server errors: %s", err)
	}
}
