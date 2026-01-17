package main

import (
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/timehub/timehub/pkg/config"
	"github.com/timehub/timehub/pkg/db"
	"github.com/timehub/timehub/pkg/logger"
	customMiddleware "github.com/timehub/timehub/pkg/middleware"

	"github.com/timehub/timehub/services/company-service/internal/delivery/http"
	"github.com/timehub/timehub/services/company-service/internal/domain"
	"github.com/timehub/timehub/services/company-service/internal/repository/postgres"
	"github.com/timehub/timehub/services/company-service/internal/usecase"

	_ "github.com/timehub/timehub/services/company-service/docs" // for swagger docs
)

// @title Company Service API
// @version 1.0
// @description Manage companies, branches, and employees.
// @host localhost:8082
// @BasePath /
func main() {
	logger.Init()
	cfg := config.Load()

	// Initialize DB (GORM)
	database, err := db.Connect(db.Config{
		Host:     "localhost",
		Port:     "5432",
		User:     "user",
		Password: "password",
		DBName:   "company_db",
		SSLMode:  "disable",
	})
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto Migrate
	err = database.AutoMigrate(
		&domain.Company{},
		&domain.Branch{},
		&domain.Category{},
		&domain.Service{},
		&domain.Employee{},
		&domain.EmployeeService{},
	)
	if err != nil {
		logger.Error("Failed to migrate database", "error", err)
		log.Fatal(err)
	}

	sqlDB, _ := database.DB()
	defer sqlDB.Close()

	logger.Info("Connected to database")

	// Init Layers
	companyRepo := postgres.NewCompanyRepository(database)
	timeout := time.Duration(2) * time.Second
	companyUsecase := usecase.NewCompanyUsecase(companyRepo, timeout)

	e := echo.New()

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(customMiddleware.RequestLogger)
	e.Use(customMiddleware.PanicRecovery)

	// Handlers
	http.NewCompanyHandler(e, companyUsecase)

	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Health Check
	e.GET("/health", func(c echo.Context) error {
		if err := sqlDB.Ping(); err != nil {
			return c.String(500, "Database unreachable")
		}
		return c.String(200, "Company Service is healthy")
	})

	logger.Info("Starting Company Service", "port", cfg.HTTPPort)
	if err := e.Start(":" + cfg.HTTPPort); err != nil {
		logger.Error("Server failed", "error", err)
	}
}
