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
	_ "github.com/vipos89/timehub/services/auth-service/docs" // for swagger docs
	"github.com/vipos89/timehub/services/auth-service/internal/delivery/http"
	"github.com/vipos89/timehub/services/auth-service/internal/domain"
	"github.com/vipos89/timehub/services/auth-service/internal/repository/postgres"
	"github.com/vipos89/timehub/services/auth-service/internal/usecase"
)

// @title Auth Service API
// @version 1.0
// @description Handle registration and login.
// @host localhost:8081
// @BasePath /
func main() {
	logger.Init()
	cfg := config.Load() // Setup config to load JWT_SECRET

	// Initialize DB (GORM) - connects to auth_db
	database, err := db.Connect(db.Config{
		Host:     "localhost",
		Port:     "5432",
		User:     "user",
		Password: "password",
		DBName:   "auth_db",
		SSLMode:  "disable",
	})
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto Migrate
	err = database.AutoMigrate(&domain.User{})
	if err != nil {
		logger.Error("Failed to migrate database", "error", err)
		log.Fatal(err)
	}

	sqlDB, _ := database.DB()
	defer sqlDB.Close()

	logger.Info("Connected to database")

	// Init Layers
	userRepo := postgres.NewUserRepository(database)
	timeout := time.Duration(2) * time.Second
	jwtSecret := "secret_key_change_me" // In real app, load from cfg.JWTSecret
	authUsecase := usecase.NewAuthUsecase(userRepo, timeout, jwtSecret)

	e := echo.New()

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(customMiddleware.RequestLogger)
	e.Use(customMiddleware.PanicRecovery)

	// Handlers
	http.NewAuthHandler(e, authUsecase)

	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Health Check
	e.GET("/health", func(c echo.Context) error {
		if err := sqlDB.Ping(); err != nil {
			return c.String(500, "Database unreachable")
		}
		return c.String(200, "Auth Service is healthy")
	})

	logger.Info("Starting Auth Service", "port", cfg.HTTPPort)
	if err := e.Start(":" + cfg.HTTPPort); err != nil {
		logger.Error("Server failed", "error", err)
		log.Fatalf("Server failed: %v", err)
	}
}
