package main

import (
	"log"
	"net/url"

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

	// Proxy Config
	authURL, _ := url.Parse(cfg.AuthServiceURL)
	companyURL, _ := url.Parse(cfg.CompanyServiceURL)
	// API Proxy
	// Auth Service
	authGroup := e.Group("/auth", middleware.Proxy(middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
		{URL: authURL},
	})))
	authGroup.Any("/*", func(c echo.Context) error { return nil })

	// Company Service
	companyGroup := e.Group("/companies", middleware.Proxy(middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
		{URL: companyURL},
	})))
	companyGroup.Any("/*", func(c echo.Context) error { return nil })

	branchGroup := e.Group("/branches", middleware.Proxy(middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
		{URL: companyURL},
	})))
	branchGroup.Any("/*", func(c echo.Context) error { return nil })

	employeeGroup := e.Group("/employees", middleware.Proxy(middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
		{URL: companyURL},
	})))
	employeeGroup.Any("/*", func(c echo.Context) error { return nil })

	servicesGroup := e.Group("/services", middleware.Proxy(middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
		{URL: companyURL},
	})))
	servicesGroup.Any("/*", func(c echo.Context) error { return nil })

	// Booking Service
	bookingURL, _ := url.Parse(cfg.BookingServiceURL)
	bookingGroup := e.Group("/bookings", middleware.Proxy(middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
		{URL: bookingURL},
	})))
	bookingGroup.Any("/*", func(c echo.Context) error { return nil })

	slotsGroup := e.Group("/slots", middleware.Proxy(middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
		{URL: bookingURL},
	})))
	slotsGroup.Any("/*", func(c echo.Context) error { return nil })

	schedulesGroup := e.Group("/schedules", middleware.Proxy(middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
		{URL: bookingURL},
	})))
	schedulesGroup.Any("/*", func(c echo.Context) error { return nil })

	shiftsGroup := e.Group("/shifts", middleware.Proxy(middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
		{URL: bookingURL},
	})))
	shiftsGroup.Any("/*", func(c echo.Context) error { return nil })

	// Swagger Proxy (Aggregate Documentation)
	// /swagger/auth/* -> auth-service/swagger/*
	e.Group("/swagger/auth", middleware.ProxyWithConfig(middleware.ProxyConfig{
		Balancer: middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
			{URL: authURL},
		}),
		Rewrite: map[string]string{
			"^/swagger/auth/*": "/swagger/$1",
		},
	}))

	// /swagger/company/* -> company-service/swagger/*
	e.Group("/swagger/company", middleware.ProxyWithConfig(middleware.ProxyConfig{
		Balancer: middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
			{URL: companyURL},
		}),
		Rewrite: map[string]string{
			"^/swagger/company/*": "/swagger/$1",
		},
	}))

	// /swagger/booking/* -> booking-service/swagger/*
	e.Group("/swagger/booking", middleware.ProxyWithConfig(middleware.ProxyConfig{
		Balancer: middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
			{URL: bookingURL},
		}),
		Rewrite: map[string]string{
			"^/swagger/booking/*": "/swagger/$1",
		},
	}))

	// Swagger
	e.File("/swagger/index.html", "docs/index_agg.html")
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
