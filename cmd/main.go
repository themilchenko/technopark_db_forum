package main

import (
        "technopark_db_forum/internal/app"
        "github.com/labstack/echo-contrib/prometheus"

        "github.com/labstack/echo/v4"
)

const (
        dbConfig = "host=84.201.157.36 port=5432 dbname=forum user=forum password=forum sslmode=disable"
)

func main() {
        e := echo.New()

        // Create new echo instance for metrics
        metricsWork := echo.New()
        p := prometheus.NewPrometheus("echo", nil)

        // Using middleware for defining wethher handler is works
        e.Use(p.HandlerFunc)
        p.SetMetricsPath(metricsWork)

        s := app.New(e)

	go func() { metricsWork.Logger.Fatal(metricsWork.Start("0.0.0.0:9091")) }()

        if err := s.Start("0.0.0.0:8080", dbConfig); err != nil {
                s.Echo.Logger.Error("server errors: %s", err)
        }
}
