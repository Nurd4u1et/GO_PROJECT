package config

import (
	"os"
)

type Config struct {
	AppPort   string
	DBUrl     string
	JWTSecret string
}

func Load() *Config {
	return &Config{
		AppPort:   getEnv("APP_PORT", "8080"),
		DBUrl:     getEnv("DATABASE_URL", "postgres://clinic_user:clinic_password@localhost:5432/clinic_db?sslmode=disable"),
		JWTSecret: getEnv("JWT_SECRET", "supersecret-dev-key"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
