package config

import (
	"os"
)

type Config struct {
	GRPCPort string
	HTTPPort string
	DBUrl    string

	// Service URLs
	AuthServiceURL    string
	CompanyServiceURL string
	BookingServiceURL string
	CRMServiceURL     string
	ReportServiceURL  string
}

func Load() *Config {
	return &Config{
		GRPCPort: getEnv("GRPC_PORT", "50051"),
		HTTPPort: getEnv("HTTP_PORT", "8080"),
		DBUrl:    getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/timehub?sslmode=disable"),

		AuthServiceURL:    getEnv("AUTH_SERVICE_URL", "http://localhost:8081"),
		CompanyServiceURL: getEnv("COMPANY_SERVICE_URL", "http://localhost:8082"),
		BookingServiceURL: getEnv("BOOKING_SERVICE_URL", "http://localhost:8083"),
		CRMServiceURL:     getEnv("CRM_SERVICE_URL", "http://localhost:8084"),
		ReportServiceURL:  getEnv("REPORT_SERVICE_URL", "http://localhost:8085"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
