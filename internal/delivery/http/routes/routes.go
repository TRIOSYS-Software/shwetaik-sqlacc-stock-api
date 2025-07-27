package routes

import (
	"shwetaik-sqlacc-stock-api/internal/delivery/http/container"
	"shwetaik-sqlacc-stock-api/internal/delivery/http/handlers"
	"shwetaik-sqlacc-stock-api/internal/delivery/http/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupRoutes(app *fiber.App, container *container.AppContainer) {

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	api := app.Group("/api/v1")

	api.Use(middleware.AuthMiddleware)
	initStockItemRoutes(api, container.StockItemHandler)
	initStockItemPriceRoutes(api, container.StockItemPriceHandler)
}

func initStockItemRoutes(api fiber.Router, handler *handlers.StockItemHandler) {

	api.Get("/stock-items", handler.GetAllStockItems)
	api.Get("/stock-items/:code", handler.GetStockItemByCode)
}

func initStockItemPriceRoutes(api fiber.Router, handler *handlers.StockItemPriceHandler) {

	api.Get("/stock-items/:code/prices", handler.GetStockItemPricesByCode)
	api.Get("/stock-items/:code/prices/:dtlKey", handler.GetStockItemPriceByDTLKey)
	api.Post("/stock-items/:code/prices", handler.CreateStockItemPrice)
	api.Put("/stock-items/:code/prices/:dtlKey", handler.UpdateStockItemPrice)
}
