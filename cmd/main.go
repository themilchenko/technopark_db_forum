package main

import (
        "technopark_db_forum/internal/app"

        "github.com/labstack/echo/v4"
)

func main() {
        e := echo.New()
        s := app.New(e)

        dbConfig := "host=db port=5432 dbname=forum user=forum password=forum sslmode=disable"

        if err := s.Start("0.0.0.0:8080", dbConfig); err != nil {
                s.Echo.Logger.Error("server errors: %s", err)
        }
}
