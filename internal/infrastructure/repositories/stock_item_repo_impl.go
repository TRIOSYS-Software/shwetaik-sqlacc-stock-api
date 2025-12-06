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
		whereClauses = append(whereClauses, "si.STOCKGROUP LIKE ?")
		args = append(args, "%"+stockGroup+"%")
	}
	if description, ok := filter["description"].(string); ok && description != "" {
		whereClauses = append(whereClauses, "si.DESCRIPTION LIKE ?")
		args = append(args, "%"+description+"%")
	}

	// Build the JOIN query
	baseQuery := `
		SELECT si.*,
			   sip.DTLKEY AS price_dtlkey, sip.CODE AS price_code, sip.TAGTYPE AS price_tagtype,
			   sip.COMPANY AS price_company, sip.SEQ AS price_seq, sip.PRICETAG AS price_pricetag,
			   sip.UOM AS price_uom, sip.QTY AS price_qty, sip.STOCKVALUE AS price_stockvalue,
			   sip.DISCOUNT AS price_discount, sip.DATEFROM AS price_datefrom, sip.DATETO AS price_dateto,
			   sip.NOTE AS price_note, sib.AUTOKEY AS barcode_autokey, sib.BARCODE AS barcode,
			   sib.UOM AS barcode_uom 
		FROM (
			SELECT si.*
			FROM ST_ITEM si
			%s
			ORDER BY si.CODE
			%s
		) si
		LEFT JOIN ST_ITEM_PRICE sip ON si.code = sip.code
		LEFT JOIN ST_ITEM_BARCODE sib ON si.code = sib.code
	`

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = " WHERE " + strings.Join(whereClauses, " AND ")
	}

	paginationSQL := ""
	if limit, ok := filter["limit"].(int); ok && limit > 0 {
		if offset, ok := filter["offset"].(int); ok && offset > 0 {
			paginationSQL = fmt.Sprintf("ROWS %d TO %d", offset+1, offset+limit)
		} else {
			paginationSQL = fmt.Sprintf("ROWS %d", limit)
		}
	}

	updateQuery := fmt.Sprintf(baseQuery, whereSQL, paginationSQL)

	// Struct for joined result
	type joinedResult struct {
		entities.STItem
		PriceDtlKey     *int     `gorm:"column:PRICE_DTLKEY"`
		PriceCode       *string  `gorm:"column:PRICE_CODE"`
		PriceTagType    *string  `gorm:"column:PRICE_TAGTYPE"`
		PriceCompany    *string  `gorm:"column:PRICE_COMPANY"`
		PriceSeq        *int     `gorm:"column:PRICE_SEQ"`
		PricePriceTag   *string  `gorm:"column:PRICE_PRICETAG"`
		PriceUOM        *string  `gorm:"column:PRICE_UOM"`
		PriceQty        *float64 `gorm:"column:PRICE_QTY"`
		PriceStockValue *float64 `gorm:"column:PRICE_STOCKVALUE"`
		PriceDiscount   *string  `gorm:"column:PRICE_DISCOUNT"`
		PriceDateFrom   *string  `gorm:"column:PRICE_DATEFROM"`
		PriceDateTo     *string  `gorm:"column:PRICE_DATETO"`
		PriceNote       *[]byte  `gorm:"column:PRICE_NOTE"`
		BarcodeAutoKey  *int     `gorm:"column:BARCODE_AUTOKEY"`
		Barcode         *string  `gorm:"column:BARCODE"`
		BarcodeUOM      *string  `gorm:"column:BARCODE_UOM"`
	}

	var joinedRows []joinedResult
	err := r.db.Raw(updateQuery, args...).Scan(&joinedRows).Error
	if err != nil {
		return nil, err
	}

	// Map results to stockItems with prices
	stockItemMap := make(map[string]*entities.STItem)
	for _, row := range joinedRows {
		item, exists := stockItemMap[row.Code]
		if !exists {
			itemCopy := row.STItem
			itemCopy.STItemPrices = []entities.STItemPrice{}
			stockItemMap[row.Code] = &itemCopy
			item = &itemCopy
		}
		if row.PriceCode != nil {
			price := entities.STItemPrice{
				DtlKey:     0,
				Code:       *row.PriceCode,
				TagType:    "",
				Company:    row.PriceCompany,
				Seq:        0,
				PriceTag:   row.PricePriceTag,
				UOM:        "",
				Qty:        row.PriceQty,
				StockValue: row.PriceStockValue,
				Discount:   row.PriceDiscount,
				Note:       row.PriceNote,
			}
			if row.PriceDtlKey != nil {
				price.DtlKey = *row.PriceDtlKey
			}
			if row.PriceTagType != nil {
				price.TagType = *row.PriceTagType
			}
			if row.PriceSeq != nil {
				price.Seq = *row.PriceSeq
			}
			if row.PriceUOM != nil {
				price.UOM = *row.PriceUOM
			}
			// DateFrom/DateTo parsing omitted for brevity
			item.STItemPrices = append(item.STItemPrices, price)
		}
		// Barcode handling can be added similarly if needed
		if row.BarcodeAutoKey != nil {
			barcode := entities.STItemBarcode{
				AutoKey: *row.BarcodeAutoKey,
				Code:    row.Code,
			}
			if row.Barcode != nil {
				barcode.Barcode = *row.Barcode
			}
			if row.BarcodeUOM != nil {
				barcode.UOM = *row.BarcodeUOM
			}
			item.STItemBarcodes = append(item.STItemBarcodes, barcode)
		}
	}

	// Convert map to slice
	stockItems := make([]entities.STItem, 0, len(stockItemMap))
	for _, item := range stockItemMap {
		stockItems = append(stockItems, *item)
	}

	return stockItems, nil
}

func (r *StockItemRepositoryImpl) GetStockItemByCode(code string) (*entities.STItem, error) {
	var stockItem entities.STItem
	err := r.db.Where("code = ?", code).Preload("STItemPrices").Preload("STItemBarcodes").First(&stockItem).Error
	return &stockItem, err
}
