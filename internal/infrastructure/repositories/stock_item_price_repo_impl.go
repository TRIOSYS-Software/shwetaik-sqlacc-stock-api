package repositories

import (
	"shwetaik-sqlacc-stock-api/internal/domain/entities"
	"shwetaik-sqlacc-stock-api/internal/domain/repositories"

	"gorm.io/gorm"
)

type StockItemPriceRepositoryImpl struct {
	db *gorm.DB
}

func NewStockItemPriceRepository(db *gorm.DB) repositories.StockItemPriceRepository {
	return &StockItemPriceRepositoryImpl{db: db}
}

func (r *StockItemPriceRepositoryImpl) GetStockItemPricesByCode(code string) ([]entities.STItemPrice, error) {
	var stockItemPrices []entities.STItemPrice
	err := r.db.Where("code = ?", code).Find(&stockItemPrices).Error
	return stockItemPrices, err
}

func (r *StockItemPriceRepositoryImpl) GetStockItemPriceByDTLKey(code string, dtlKey int) (*entities.STItemPrice, error) {
	var stockItemPrice entities.STItemPrice
	err := r.db.Where("code = ?", code).Where("dtlkey = ?", dtlKey).First(&stockItemPrice).Error
	return &stockItemPrice, err
}

