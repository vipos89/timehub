package config

import (
	"os"
)

type Config struct {
	GRPCPort string
	HTTPPort string
	DBUrl    string
}

func Load() *Config {
	return &Config{
		GRPCPort: getEnv("GRPC_PORT", "50051"),
		HTTPPort: getEnv("HTTP_PORT", "8080"),
		DBUrl:    getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/timehub?sslmode=disable"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
