package config

import (
	"os"

	fiberlog "github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

var Cfg *Config

type Config struct {
	Host       string
	Port       string
	DBString   string
	JWT_SECRET string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		fiberlog.Errorf("Error loading .env file: %v", err)
	}

	var config Config
	config.Host = getEnv("HOST", "localhost")
	config.Port = getEnv("PORT", "8080")
	config.DBString = getEnv("DB_STRING", "")
	config.JWT_SECRET = getEnv("JWT_SECRET", "")

	return &config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	if defaultValue == "" {
		fiberlog.Fatalf("environment variable %q not set", key)
	}
	return defaultValue
}

func init() {
	Cfg = Load()
}
