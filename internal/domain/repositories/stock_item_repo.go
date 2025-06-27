package repositories

import (
	"shwetaik-sqlacc-stock-api/internal/domain/entities"
)

type StockItemRepository interface {
	GetAllStockItems(filter map[string]any) ([]entities.STItem, error)
	GetStockItemByCode(code string) (*entities.STItem, error)
}
