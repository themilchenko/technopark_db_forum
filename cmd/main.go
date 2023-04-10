package main

import (
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"os"

	"technopark_db_forum/internal/app"
)

const (
	dbConfig = "host=postgresql port=5432 user=forum dbname=forum password=forum sslmode=disable"
)

func main() {
 e := echo.New()

	// Create new echo instance for metrics
	metricsWork := echo.New()
	p := prometheus.NewPrometheus("echo", nil)

	// // Using middleware for defining wethher handler is works
	e.Use(p.HandlerFunc)
	p.SetMetricsPath(metricsWork)

	s := app.New(e)

	go func() { metricsWork.Logger.Fatal(metricsWork.Start(":" + os.Getenv("METRICS_PORT"))) }()

	if err := s.Start(":" + os.Getenv("PORT"), dbConfig); err != nil {
		s.Echo.Logger.Error("server errors: %s", err )
	}
}
