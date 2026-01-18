package main

import (
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/vipos89/timehub/pkg/config"
	"github.com/vipos89/timehub/pkg/db"
	"github.com/vipos89/timehub/pkg/logger"
	customMiddleware "github.com/vipos89/timehub/pkg/middleware"

	_ "github.com/vipos89/timehub/services/booking-service/docs" // for swagger docs
	"github.com/vipos89/timehub/services/booking-service/internal/delivery/http"
	"github.com/vipos89/timehub/services/booking-service/internal/domain"
	"github.com/vipos89/timehub/services/booking-service/internal/repository/postgres"
	"github.com/vipos89/timehub/services/booking-service/internal/usecase"
)

// @title Booking Service API
// @version 1.0
// @description Manage schedules and appointments.
// @host localhost:8083
// @BasePath /
func main() {
	logger.Init()
	cfg := config.Load()

	// Initialize DB (GORM) - connects to booking_db
	database, err := db.ConnectDSN(cfg.DBUrl)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto Migrate
	err = database.AutoMigrate(
		&domain.Schedule{},
		&domain.WorkShift{},
		&domain.Appointment{},
	)
	if err != nil {
		logger.Error("Failed to migrate database", "error", err)
		log.Fatal(err)
	}

	sqlDB, _ := database.DB()
	defer sqlDB.Close()

	logger.Info("Connected to database")

	// Init Layers
	bookingRepo := postgres.NewBookingRepository(database)
	timeout := time.Duration(2) * time.Second
	bookingUsecase := usecase.NewBookingUsecase(bookingRepo, timeout)

	e := echo.New()

	// Middleware
	e.Use(middleware.Recover())
	e.Use(customMiddleware.RequestLogger)
	e.Use(customMiddleware.PanicRecovery)

	// Handlers
	http.NewBookingHandler(e, bookingUsecase)

	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Health Check
	e.GET("/health", func(c echo.Context) error {
		if err := sqlDB.Ping(); err != nil {
			return c.String(500, "Database unreachable")
		}
		return c.String(200, "Booking Service is healthy")
	})

	logger.Info("Starting Booking Service", "port", cfg.HTTPPort)
	if err := e.Start(":" + cfg.HTTPPort); err != nil {
		logger.Error("Server failed", "error", err)
		log.Fatalf("Server failed: %v", err)
	}
}
