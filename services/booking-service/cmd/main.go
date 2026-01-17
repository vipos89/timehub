package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/timehub/timehub/pkg/config"
	"github.com/timehub/timehub/pkg/logger"
	_ "github.com/timehub/timehub/services/booking-service/docs" // for swagger docs
)

// @title Booking Service API
// @version 1.0
// @description Manage appointments and availability.
// @host localhost:8083
// @BasePath /
func main() {
	logger.Init()
	cfg := config.Load()

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Health Check
	e.GET("/health", func(c echo.Context) error {
		return c.String(200, "Booking Service is health")
	})

	logger.Info("Starting Booking Service", "port", cfg.HTTPPort)
	if err := e.Start(":" + cfg.HTTPPort); err != nil {
		logger.Error("Server failed", "error", err)
		log.Fatalf("Server failed: %v", err)
	}
}
