package repositories

import (
	"shwetaik-sqlacc-stock-api/internal/domain/entities"
)

type StockItemRepository interface {
	GetAllStockItems(filter map[string]any) ([]entities.STItem, error)
	// GetStockItemByCode returns the stock item's balance for the given
	// location, or summed across all locations when location is "".
	GetStockItemByCode(code string, location string) (*entities.STItem, error)
}
