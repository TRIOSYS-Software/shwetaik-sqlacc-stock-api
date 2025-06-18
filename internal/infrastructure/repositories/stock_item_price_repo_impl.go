package repositories

import (
	"shwetaik-sqlacc-stock-api/internal/domain/entities"
	"shwetaik-sqlacc-stock-api/internal/domain/repositories"
	"time"

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

func (r *StockItemPriceRepositoryImpl) CreateStockItemPrice(code string, stockItemPrice *entities.STItemPrice) error {
	// var lastId int
	// err := r.db.Model(&entities.STItemPrice{}).Select("MAX(dtlkey)").Scan(&lastId).Error
	// if err != nil {
	// 	return err
	// }
	// stockItemPrice.DtlKey = lastId + 1
	// return r.db.Where("code = ?", code).Create(stockItemPrice).Error
	tx := r.db.Begin()
	var lastId int
	err := tx.Model(&entities.STItemPrice{}).Select("MIN(dtlkey)").Scan(&lastId).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	if lastId > 0 {
		stockItemPrice.DtlKey = -1
	} else {
		stockItemPrice.DtlKey = lastId - 1
	}
	var maxSeq int
	err = tx.Model(&entities.STItemPrice{}).Select("MAX(seq)").Where("code = ?", code).Scan(&maxSeq).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	stockItemPrice.Seq = maxSeq + 1000
	if err := tx.Where("code = ?", code).Create(stockItemPrice).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&entities.STItem{}).Where("code = ?", code).Update("LASTMODIFIED", time.Now().Unix()).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *StockItemPriceRepositoryImpl) UpdateStockItemPrice(code string, stockItemPrice *entities.STItemPrice) error {
	return r.db.Where("code = ?", code).Save(stockItemPrice).Error
}
