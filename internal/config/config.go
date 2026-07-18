package config

import (
	"os"
	"strings"

	fiberlog "github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

var Cfg *Config

type Config struct {
	Host        string
	Port        string
	DBString    string
	JWT_SECRET  string
	ServiceName string

	SQLAccountAPIHost         string
	SQLAccountAPIAccessKey    string
	SQLAccountAPISecretKey    string
	SQLAccountAPIRegion       string
	SQLAccountAPIService      string
	SQLAccountAPISessionToken string

	WebhookURLs []string
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
	config.ServiceName = getEnv("SERVICE_NAME", "stock-api")

	// Vendor (SQL Account) API credentials — optional until the wrapped
	// write endpoints are actually used, so these don't fail startup.
	config.SQLAccountAPIHost = os.Getenv("SQLACCOUNT_API_HOST")
	config.SQLAccountAPIAccessKey = os.Getenv("SQLACCOUNT_API_ACCESS_KEY")
	config.SQLAccountAPISecretKey = os.Getenv("SQLACCOUNT_API_SECRET_KEY")
	config.SQLAccountAPIRegion = os.Getenv("SQLACCOUNT_API_REGION")
	config.SQLAccountAPIService = os.Getenv("SQLACCOUNT_API_SERVICE")
	config.SQLAccountAPISessionToken = os.Getenv("SQLACCOUNT_API_SESSION_TOKEN")

	// Optional, comma-separated — the stock item monitor no-ops on webhook
	// delivery if unset.
	config.WebhookURLs = parseWebhookURLs(os.Getenv("WEBHOOK_URL"))

	return &config
}

func parseWebhookURLs(raw string) []string {
	if raw == "" {
		return nil
	}
	var urls []string
	for url := range strings.SplitSeq(raw, ",") {
		if url = strings.TrimSpace(url); url != "" {
			urls = append(urls, url)
		}
	}
	return urls
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
