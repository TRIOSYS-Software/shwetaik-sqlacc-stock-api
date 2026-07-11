package repositories

import (
	"fmt"
	"shwetaik-sqlacc-stock-api/internal/domain/entities"
	"shwetaik-sqlacc-stock-api/internal/domain/repositories"
	"strings"

	"gorm.io/gorm"
)

type StockItemRepositoryImpl struct {
	db *gorm.DB
}

func NewStockItemRepository(db *gorm.DB) repositories.StockItemRepository {
	return &StockItemRepositoryImpl{db: db}
}

func (r *StockItemRepositoryImpl) GetAllStockItems(filter map[string]any) ([]entities.STItem, error) {
	// Build WHERE clauses and args
	whereClauses := []string{}
	args := []interface{}{}

	if stockGroup, ok := filter["stock_group"].(string); ok && stockGroup != "" {
		whereClauses = append(whereClauses, "STOCKGROUP LIKE ?")
		args = append(args, "%"+stockGroup+"%")
	}
	if description, ok := filter["description"].(string); ok && description != "" {
		whereClauses = append(whereClauses, "DESCRIPTION LIKE ?")
		args = append(args, "%"+description+"%")
	}
	// Keyset pagination: the caller passes the CODE of the last item from
	// the previous page instead of an offset, so the DB only has to find
	// rows past that point rather than re-materializing every prior page.
	if after, ok := filter["after"].(string); ok && after != "" {
		whereClauses = append(whereClauses, "CODE > ?")
		args = append(args, after)
	}

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = " WHERE " + strings.Join(whereClauses, " AND ")
	}

	limit, hasLimit := filter["limit"].(int)
	if !hasLimit || limit <= 0 {
		limit = 0
	}

	paginationSQL := ""
	if limit > 0 {
		paginationSQL = fmt.Sprintf("ROWS %d", limit)
	}

	// Reused as a derived table below so prices/barcodes are joined against
	// the same filtered/paginated code set in SQL, rather than collected
	// into a Go slice and matched with a large IN (...) list — Firebird's
	// optimizer handles a join on an indexed CODE column far better than an
	// IN list of hundreds/thousands of literal values, which tends to fall
	// back to a per-row scan instead of index seeks.
	filteredCodesSQL := fmt.Sprintf(`
		SELECT CODE FROM ST_ITEM
		%s
		ORDER BY CODE
		%s
	`, whereSQL, paginationSQL)

	// Only select columns the response DTO actually uses — ST_ITEM has
	// several BLOB columns (DESCRIPTION3, PICTURE, ATTACHMENTS, NOTE) that
	// are otherwise fetched and discarded for every row, and Firebird reads
	// each non-null BLOB as a separate round trip.
	itemQuery := fmt.Sprintf(`
		SELECT DOCKEY, CODE, DESCRIPTION, STOCKGROUP
		FROM ST_ITEM
		%s
		ORDER BY CODE
		%s
	`, whereSQL, paginationSQL)

	// Fetch the filtered/paginated items first; this query alone determines
	// the item set and its order.
	var stockItems []entities.STItem
	if err := r.db.Raw(itemQuery, args...).Scan(&stockItems).Error; err != nil {
		return nil, err
	}
	if len(stockItems) == 0 {
		return stockItems, nil
	}

	itemIndexByCode := make(map[string]int, len(stockItems))
	for i, item := range stockItems {
		itemIndexByCode[item.Code] = i
		stockItems[i].STItemPrices = []entities.STItemPrice{}
		stockItems[i].STItemBarcodes = []entities.STItemBarcode{}
	}

	// Prices and barcodes are each their own one-to-many relation to ST_ITEM,
	// so they're joined independently and merged by code in Go — joining
	// both in a single query would cross-multiply rows (fan-out duplication).
	priceQuery := fmt.Sprintf(`
		SELECT sip.DTLKEY, sip.CODE, sip.PRICETAG, sip.UOM, sip.QTY, sip.STOCKVALUE
		FROM (%s) fc
		JOIN ST_ITEM_PRICE sip ON fc.CODE = sip.CODE
	`, filteredCodesSQL)
	var prices []entities.STItemPrice
	if err := r.db.Raw(priceQuery, args...).Scan(&prices).Error; err != nil {
		return nil, err
	}
	for _, price := range prices {
		if idx, ok := itemIndexByCode[price.Code]; ok {
			stockItems[idx].STItemPrices = append(stockItems[idx].STItemPrices, price)
		}
	}

	barcodeQuery := fmt.Sprintf(`
		SELECT sib.AUTOKEY, sib.CODE, sib.BARCODE, sib.UOM
		FROM (%s) fc
		JOIN ST_ITEM_BARCODE sib ON fc.CODE = sib.CODE
	`, filteredCodesSQL)
	var barcodes []entities.STItemBarcode
	if err := r.db.Raw(barcodeQuery, args...).Scan(&barcodes).Error; err != nil {
		return nil, err
	}
	for _, barcode := range barcodes {
		if idx, ok := itemIndexByCode[barcode.Code]; ok {
			stockItems[idx].STItemBarcodes = append(stockItems[idx].STItemBarcodes, barcode)
		}
	}

	return stockItems, nil
}

func (r *StockItemRepositoryImpl) GetStockItemByCode(code string) (*entities.STItem, error) {
	var stockItem entities.STItem
	err := r.db.Where("code = ?", code).Preload("STItemPrices").Preload("STItemBarcodes").First(&stockItem).Error
	return &stockItem, err
}
