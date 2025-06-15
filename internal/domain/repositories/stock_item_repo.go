package repositories

import (
	"shwetaik-sqlacc-stock-api/internal/domain/entities"
)

type StockItemRepository interface {
	GetAllStockItems(limit int, offset int) ([]entities.STItem, error)
	GetStockItemByCode(code string) (*entities.STItem, error)
}
