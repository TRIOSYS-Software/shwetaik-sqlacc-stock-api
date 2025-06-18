package repositories

import (
	"shwetaik-sqlacc-stock-api/internal/domain/entities"
	"shwetaik-sqlacc-stock-api/internal/domain/repositories"

	"gorm.io/gorm"
)

type StockItemRepositoryImpl struct {
	db *gorm.DB
}

func NewStockItemRepository(db *gorm.DB) repositories.StockItemRepository {
	return &StockItemRepositoryImpl{db: db}
}

func (r *StockItemRepositoryImpl) GetAllStockItems(limt int, offset int) ([]entities.STItem, error) {
	var stockItems []entities.STItem
	query := r.db.Model(&entities.STItem{}).Preload("STItemPrices")
	if limt > 0 {
		query = query.Limit(limt)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	err := query.Find(&stockItems).Error
	return stockItems, err
}

func (r *StockItemRepositoryImpl) GetStockItemByCode(code string) (*entities.STItem, error) {
	var stockItem entities.STItem
	err := r.db.Where("code = ?", code).Preload("STItemPrices").First(&stockItem).Error
	return &stockItem, err
}
