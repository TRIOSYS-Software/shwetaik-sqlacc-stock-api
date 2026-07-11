package dto

type StockItemPriceResponse struct {
	DtlKey     int      `json:"dtl_key"`
	Code       string   `json:"code"`
	PriceTag   *string  `json:"price_tag"`
	UOM        string   `json:"uom"`
	Qty        *float64 `json:"qty"`
	StockValue *float64 `json:"stock_value"`
}

type StockItemBarcodeResponse struct {
	AutoKey int    `json:"autokey"`
	Barcode string `json:"barcode"`
	UOM     string `json:"uom"`
}

// StockItemPriceItem is one line in a PutStockItemPricesRequest. DtlKey
// distinguishes create from update: nil (field omitted) means create a new
// price line, a set value means replace the existing line with that dtlkey.
type StockItemPriceItem struct {
	DtlKey     *int     `json:"dtlkey,omitempty"`
	PriceTag   string   `json:"price_tag"`
	UOM        string   `json:"uom"`
	Qty        *float64 `json:"qty"`
	StockValue *float64 `json:"stock_value"`
}

// PutStockItemPricesRequest replaces a stock item's entire customer price
// list — this is a true replace, not a merge, so it must include every
// price line the caller wants kept.
type PutStockItemPricesRequest struct {
	Prices []StockItemPriceItem `json:"prices"`
}
