package main

import (
	"flag"
	"shwetaik-sqlacc-stock-api/internal/config"
	"shwetaik-sqlacc-stock-api/internal/delivery/http/container"
	"shwetaik-sqlacc-stock-api/internal/delivery/http/routes"
	"shwetaik-sqlacc-stock-api/internal/infrastructure/database"
	"shwetaik-sqlacc-stock-api/internal/infrastructure/monitor"
	"shwetaik-sqlacc-stock-api/internal/infrastructure/webhook"
	"shwetaik-sqlacc-stock-api/scripts"
	"time"

	"github.com/gofiber/fiber/v2"
	fiberlog "github.com/gofiber/fiber/v2/log"
)

func main() {
	generateToken := flag.Bool("generate-token", false, "Generate API Token")
	flag.Parse()
	if *generateToken {
		scripts.GenerateJWTToken()
		return
	}

	cfg := config.Load()

	db, err := database.NewConnection(cfg)
	if err != nil {
		fiberlog.Fatalf("Error connecting to database: %v", err)
	}

	webhookClient := webhook.NewClient(cfg.WebhookURLs)
	monitor.StartStockItemChangeMonitor(db, webhookClient, 30*time.Second)

	app := fiber.New()

	container := container.NewAppContainer(db, cfg)
	routes.SetupRoutes(app, container)
	if err := app.Listen(cfg.Host + ":" + cfg.Port); err != nil {
		fiberlog.Fatalf("Error starting server: %v", err)
	}

}
