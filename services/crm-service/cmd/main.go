package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/vipos89/timehub/pkg/config"
	"github.com/vipos89/timehub/pkg/logger"
	_ "github.com/vipos89/timehub/services/crm-service/docs" // for swagger docs
)

// @title CRM Service API
// @version 1.0
// @description Manage client base and history.
// @host localhost:8084
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
		return c.String(200, "CRM Service is health")
	})

	logger.Info("Starting CRM Service", "port", cfg.HTTPPort)
	if err := e.Start(":" + cfg.HTTPPort); err != nil {
		logger.Error("Server failed", "error", err)
		log.Fatalf("Server failed: %v", err)
	}
}
