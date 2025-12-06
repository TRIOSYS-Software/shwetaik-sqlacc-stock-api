package container

import (
	"shwetaik-sqlacc-stock-api/internal/delivery/http/handlers"
	"shwetaik-sqlacc-stock-api/internal/infrastructure/repositories"

	"shwetaik-sqlacc-stock-api/internal/usecases"

	"gorm.io/gorm"
)

type AppContainer struct {
	StockItemHandler      *handlers.StockItemHandler
	StockItemPriceHandler *handlers.StockItemPriceHandler
}

func NewAppContainer(db *gorm.DB) *AppContainer {
	stockItemRepo := repositories.NewStockItemRepository(db)
	stockItemUsecase := usecases.NewStockItemUseCase(stockItemRepo)
	stockItemHandler := handlers.NewStockItemHandler(stockItemUsecase)

	stockItemPriceRepo := repositories.NewStockItemPriceRepository(db)
	stockItemPriceUsecase := usecases.NewStockItemPriceUseCase(stockItemPriceRepo)
	stockItemPriceHandler := handlers.NewStockItemPriceHandler(stockItemPriceUsecase)

	return &AppContainer{
		StockItemHandler:      stockItemHandler,
		StockItemPriceHandler: stockItemPriceHandler,
	}
}
