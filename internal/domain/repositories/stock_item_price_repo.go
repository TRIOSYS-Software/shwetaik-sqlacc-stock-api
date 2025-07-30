package repositories

import (
	"shwetaik-sqlacc-stock-api/internal/domain/entities"
)

type StockItemPriceRepository interface {
	GetStockItemPricesByCode(code string) ([]entities.STItemPrice, error)
	GetStockItemPriceByDTLKey(code string, dtlKey int) (*entities.STItemPrice, error)
	CreateStockItemPrice(stockItemPrice *entities.STItemPrice) error
	UpdateStockItemPrice(code string, stockItemPrice *entities.STItemPrice) error
}
