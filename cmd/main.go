package main

import (
	"shwetaik-sqlacc-stock-api/internal/config"
	"shwetaik-sqlacc-stock-api/internal/delivery/http/container"
	"shwetaik-sqlacc-stock-api/internal/delivery/http/routes"
	"shwetaik-sqlacc-stock-api/internal/infrastructure/database"

	"github.com/gofiber/fiber/v2"
	fiberlog "github.com/gofiber/fiber/v2/log"
)

func main() {
	cfg := config.Load()

	db, err := database.NewConnection(cfg)
	if err != nil {
		fiberlog.Fatalf("Error connecting to database: %v", err)
	}

	app := fiber.New()

	container := container.NewAppContainer(db)
	routes.SetupRoutes(app, container)
	if err := app.Listen(cfg.Host + ":" + cfg.Port); err != nil {
		fiberlog.Fatalf("Error starting server: %v", err)
	}

}
