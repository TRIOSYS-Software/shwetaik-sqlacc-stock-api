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
	// When a location is given, restrict to items that actually have stock
	// there — this has to be a WHERE condition (not a post-fetch filter),
	// otherwise pagination/"after" would be computed against the full
	// unfiltered item set and pages could come back short or skip items.
	location, _ := filter["location"].(string)
	if location != "" {
		whereClauses = append(whereClauses, "EXISTS (SELECT 1 FROM ST_TR tr WHERE tr.ITEMCODE = ST_ITEM.CODE AND tr.LOCATION = ?)")
		args = append(args, location)
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
		if location != "" {
			loc := location
			stockItems[i].Location = &loc
		}
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

	// Balance comes from ST_TR (stock transactions), joined against the same
	// filtered/paginated code set rather than an IN (...) list of codes —
	// ST_TR is a transaction ledger, likely much larger than ST_ITEM, so the
	// same index-seek-over-IN-list reasoning as prices/barcodes applies.
	// Items with no matching ST_TR rows simply keep the zero-value Balance.
	balanceWhereSQL := ""
	if location != "" {
		balanceWhereSQL = "WHERE tr.LOCATION = ?"
	}
	balanceQuery := fmt.Sprintf(`
		SELECT tr.ITEMCODE, SUM(tr.QTY) AS CURRENT_BALANCE
		FROM (%s) fc
		JOIN ST_TR tr ON fc.CODE = tr.ITEMCODE
		%s
		GROUP BY tr.ITEMCODE
	`, filteredCodesSQL, balanceWhereSQL)

	balanceArgs := append([]interface{}{}, args...)
	if location != "" {
		balanceArgs = append(balanceArgs, location)
	}

	type stockItemBalance struct {
		ItemCode string  `gorm:"column:ITEMCODE"`
		Balance  float64 `gorm:"column:CURRENT_BALANCE"`
	}
	var balances []stockItemBalance
	if err := r.db.Raw(balanceQuery, balanceArgs...).Scan(&balances).Error; err != nil {
		return nil, err
	}
	for _, b := range balances {
		if idx, ok := itemIndexByCode[b.ItemCode]; ok {
			stockItems[idx].Balance = b.Balance
		}
	}

	return stockItems, nil
}

func (r *StockItemRepositoryImpl) GetStockItemByCode(code string, location string) (*entities.STItem, error) {
	// Select only the columns the response DTO uses and fetch prices and
	// barcodes as separate targeted queries rather than GORM Preload, which
	// otherwise runs SELECT * on both — pulling ST_ITEM's BLOB columns
	// (PICTURE, ATTACHMENTS, NOTE, DESCRIPTION3) and every price row's NOTE
	// BLOB, none of which the response ever uses.
	var stockItem entities.STItem
	if err := r.db.Select("DOCKEY", "CODE", "DESCRIPTION", "STOCKGROUP").
		Where("code = ?", code).First(&stockItem).Error; err != nil {
		return nil, err
	}

	var prices []entities.STItemPrice
	if err := r.db.Select("DTLKEY", "CODE", "PRICETAG", "UOM", "QTY", "STOCKVALUE").
		Where("code = ?", code).Find(&prices).Error; err != nil {
		return nil, err
	}
	stockItem.STItemPrices = prices

	var barcodes []entities.STItemBarcode
	if err := r.db.Where("code = ?", code).Find(&barcodes).Error; err != nil {
		return nil, err
	}
	stockItem.STItemBarcodes = barcodes

	balanceQuery := "SELECT COALESCE(SUM(QTY), 0) FROM ST_TR WHERE ITEMCODE = ?"
	balanceArgs := []interface{}{code}
	if location != "" {
		balanceQuery += " AND LOCATION = ?"
		balanceArgs = append(balanceArgs, location)
		stockItem.Location = &location
	}
	if err := r.db.Raw(balanceQuery, balanceArgs...).Scan(&stockItem.Balance).Error; err != nil {
		return nil, err
	}

	return &stockItem, nil
}
