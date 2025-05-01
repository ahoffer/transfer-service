package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     int
	LogLevel string
	DBConfig DBConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func Load() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using environment variables or defaults")
	}

	portStr := getEnv("PORT", "8080")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 8080 // Fallback to default port if parsing fails
	}

	return &Config{
		Port:     port,
		LogLevel: getEnv("LOG_LEVEL", "info"),
		DBConfig: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "transfer_service"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
