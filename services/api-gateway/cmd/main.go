package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/vipos89/timehub/pkg/config"
	"github.com/vipos89/timehub/pkg/logger"
	_ "github.com/vipos89/timehub/services/api-gateway/docs" // for swagger docs
)

// @title TimeHub API Gateway
// @version 1.0
// @description Entry point for TimeHub microservices.
// @host localhost:8080
// @BasePath /
func main() {
	logger.Init()
	cfg := config.Load()

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Health Check
	// @Summary Health Check
	// @Description Check if the service is running
	// @Tags health
	// @Produce plain
	// @Success 200 {string} string "API Gateway is running"
	// @Router /health [get]
	e.GET("/health", func(c echo.Context) error {
		return c.String(200, "API Gateway is running")
	})

	logger.Info("Starting API Gateway", "port", cfg.HTTPPort)
	if err := e.Start(":" + cfg.HTTPPort); err != nil {
		logger.Error("Failed to start server", "error", err)
		log.Fatalf("Failed to start server: %v", err)
	}
}
