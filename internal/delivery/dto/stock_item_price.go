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

type StockItemPriceRequest struct {
	// Code       string   `json:"code"`
	PriceTag   string   `json:"price_tag"`
	UOM        string   `json:"uom"`
	Qty        *float64 `json:"qty"`
	StockValue *float64 `json:"stock_value"`
}

type BulkUpdateStockItemPriceRequest struct {
	Code       string   `json:"code"`
	PriceTag   string   `json:"price_tag"`
	UOM        string   `json:"uom"`
	Qty        *float64 `json:"qty"`
	StockValue *float64 `json:"stock_value"`
}
