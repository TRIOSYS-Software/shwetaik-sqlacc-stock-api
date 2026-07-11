package container

import (
	"shwetaik-sqlacc-stock-api/internal/config"
	"shwetaik-sqlacc-stock-api/internal/delivery/http/handlers"
	"shwetaik-sqlacc-stock-api/internal/infrastructure/repositories"
	"shwetaik-sqlacc-stock-api/internal/infrastructure/sqlaccountapi"

	"shwetaik-sqlacc-stock-api/internal/usecases"

	"gorm.io/gorm"
)

type AppContainer struct {
	StockItemHandler      *handlers.StockItemHandler
	StockItemPriceHandler *handlers.StockItemPriceHandler
	GLAccHandler          *handlers.GLAccHandler
	PaymentMethodHandler  *handlers.PaymentMethodHandler
	ProjectHandler        *handlers.ProjectHandler
	PaymentVoucherHandler *handlers.PaymentVoucherHandler
}

func NewAppContainer(db *gorm.DB, cfg *config.Config) *AppContainer {
	vendorAPIClient := sqlaccountapi.NewClient(cfg)

	stockItemRepo := repositories.NewStockItemRepository(db)
	stockItemUsecase := usecases.NewStockItemUseCase(stockItemRepo)
	stockItemHandler := handlers.NewStockItemHandler(stockItemUsecase)

	stockItemPriceRepo := repositories.NewStockItemPriceRepository(db)
	stockItemPriceUsecase := usecases.NewStockItemPriceUseCase(stockItemPriceRepo, stockItemRepo, vendorAPIClient)
	stockItemPriceHandler := handlers.NewStockItemPriceHandler(stockItemPriceUsecase)

	glAccRepo := repositories.NewGLAccRepository(db)
	glAccUsecase := usecases.NewGLAccUseCase(glAccRepo)
	glAccHandler := handlers.NewGLAccHandler(glAccUsecase)

	paymentMethodRepo := repositories.NewPaymentMethodRepository(db)
	paymentMethodUsecase := usecases.NewPaymentMethodUseCase(paymentMethodRepo)
	paymentMethodHandler := handlers.NewPaymentMethodHandler(paymentMethodUsecase)

	projectRepo := repositories.NewProjectRepository(db)
	projectUsecase := usecases.NewProjectUseCase(projectRepo)
	projectHandler := handlers.NewProjectHandler(projectUsecase)

	paymentVoucherUsecase := usecases.NewPaymentVoucherUseCase(vendorAPIClient)
	paymentVoucherHandler := handlers.NewPaymentVoucherHandler(paymentVoucherUsecase)

	return &AppContainer{
		StockItemHandler:      stockItemHandler,
		StockItemPriceHandler: stockItemPriceHandler,
		GLAccHandler:          glAccHandler,
		PaymentMethodHandler:  paymentMethodHandler,
		ProjectHandler:        projectHandler,
		PaymentVoucherHandler: paymentVoucherHandler,
	}
}
