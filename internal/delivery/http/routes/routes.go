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
	initGLAccRoutes(api, container.GLAccHandler)
	initPaymentMethodRoutes(api, container.PaymentMethodHandler)
	initProjectRoutes(api, container.ProjectHandler)
	initPaymentVoucherRoutes(api, container.PaymentVoucherHandler)
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

func initGLAccRoutes(api fiber.Router, handler *handlers.GLAccHandler) {

	api.Get("/gl-accounts", handler.GetAllGLAccs)
	api.Get("/gl-accounts/:code", handler.GetGLAccByCode)
}

func initPaymentMethodRoutes(api fiber.Router, handler *handlers.PaymentMethodHandler) {

	api.Get("/payment-methods", handler.GetAllPaymentMethods)
}

func initProjectRoutes(api fiber.Router, handler *handlers.ProjectHandler) {

	api.Get("/projects", handler.GetAllProjects)
	api.Get("/projects/:code", handler.GetProjectByCode)
}

func initPaymentVoucherRoutes(api fiber.Router, handler *handlers.PaymentVoucherHandler) {

	api.Post("/payment-vouchers", handler.CreatePaymentVoucher)
}
