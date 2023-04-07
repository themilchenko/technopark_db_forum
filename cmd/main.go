package main

import (
	// "github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"

	"technopark_db_forum/internal/app"
)

const (
	dbConfig = "host=localhost port=5432 user=postgres dbname=postgres sslmode=disable"
)

func main() {
	e := echo.New()

	// Create new echo instance for metrics
	// metricsWork := echo.New()
	// p := prometheus.NewPrometheus("echo", nil)

	// // Using middleware for defining wethher handler is works
	// e.Use(p.HandlerFunc)
	// p.SetMetricsPath(metricsWork)

	s := app.New(e)

	// go func() { metricsWork.Logger.Fatal(metricsWork.Start("0.0.0.0:9091")) }()

	if err := s.Start(":8080", dbConfig); err != nil {
		s.Echo.Logger.Error("server errors: %s", err )
	}
}
