package routes

import (
	"shwetaik-sqlacc-stock-api/internal/delivery/http/handlers"
	"shwetaik-sqlacc-stock-api/internal/delivery/http/middleware"
	"shwetaik-sqlacc-stock-api/internal/infrastructure/repositories"
	"shwetaik-sqlacc-stock-api/internal/usecases"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	api := app.Group("/api/v1")

	api.Use(middleware.AuthMiddleware)
	initStockItemRoutes(api, db)
}

func initStockItemRoutes(api fiber.Router, db *gorm.DB) {
	stockItemRepo := repositories.NewStockItemRepository(db)
	stockItemUsecase := usecases.NewStockItemUseCase(stockItemRepo)
	stockItemHandler := handlers.NewStockItemHandler(stockItemUsecase)

	api.Get("/stock-items", stockItemHandler.GetAllStockItems)
}
