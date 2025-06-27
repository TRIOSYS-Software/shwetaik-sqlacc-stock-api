package repositories

import (
	"shwetaik-sqlacc-stock-api/internal/domain/entities"
	"shwetaik-sqlacc-stock-api/internal/domain/repositories"
	"shwetaik-sqlacc-stock-api/pkg/utils"
	"sync"

	"gorm.io/gorm"
)

type StockItemRepositoryImpl struct {
	db *gorm.DB
}

func NewStockItemRepository(db *gorm.DB) repositories.StockItemRepository {
	return &StockItemRepositoryImpl{db: db}
}

func (r *StockItemRepositoryImpl) GetAllStockItems(filter map[string]any) ([]entities.STItem, error) {
	var stockItems []entities.STItem
	query := r.db.Model(&entities.STItem{})

	if limit, ok := filter["limit"].(int); ok && limit > 0 {
		query = query.Limit(limit)
		if offset, ok := filter["offset"].(int); ok && offset > 0 {
			query = query.Offset(offset)
		}
	}
	if stockGroup, ok := filter["stock_group"].(string); ok && stockGroup != "" {
		query = query.Where("STOCKGROUP LIKE ?", "%"+stockGroup+"%")
	}
	if description, ok := filter["description"].(string); ok && description != "" {
		query = query.Where("DESCRIPTION LIKE ?", "%"+description+"%")
	}

	err := query.Find(&stockItems).Error

	codes := make([]string, len(stockItems))
	for i, stockItem := range stockItems {
		codes[i] = stockItem.Code
	}

	chunked := utils.ChunkSlice(codes, 1500)
	priceMap := sync.Map{}
	var wg sync.WaitGroup
	errChan := make(chan error, len(chunked))
	semaphore := make(chan bool, 5)
	for _, chunk := range chunked {
		wg.Add(1)
		go func() {
			defer wg.Done()
			semaphore <- true
			defer func() { <-semaphore }()
			var stockItemPrices []entities.STItemPrice
			err := r.db.Where("code IN ?", chunk).Find(&stockItemPrices).Error
			if err != nil {
				errChan <- err
				return
			}
			for _, stockItemPrice := range stockItemPrices {
				val, _ := priceMap.LoadOrStore(stockItemPrice.Code, []entities.STItemPrice{})
				priceMap.Store(stockItemPrice.Code, append(val.([]entities.STItemPrice), stockItemPrice))
			}
		}()
	}

	wg.Wait()
	close(errChan)
	if len(errChan) > 0 {
		return nil, <-errChan
	}

	for i, item := range stockItems {
		if prices, ok := priceMap.Load(item.Code); ok {
			stockItems[i].STItemPrices = prices.([]entities.STItemPrice)
		}
	}

	return stockItems, err
}

func (r *StockItemRepositoryImpl) GetStockItemByCode(code string) (*entities.STItem, error) {
	var stockItem entities.STItem
	err := r.db.Where("code = ?", code).Preload("STItemPrices").First(&stockItem).Error
	return &stockItem, err
}
