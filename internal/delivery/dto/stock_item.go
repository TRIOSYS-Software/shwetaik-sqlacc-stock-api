package dto

type StockItemResponse struct {
	DocKey         int                        `json:"dockey"`
	Code           string                     `json:"code"`
	Description    string                     `json:"description"`
	StockGroup     string                     `json:"stock_group"`
	Balance        float64                    `json:"balance"`
	Location       string                     `json:"location,omitempty"`
	STItemPrices   []StockItemPriceResponse   `json:"stock_item_prices"`
	STItemBarcodes []StockItemBarcodeResponse `json:"stock_item_barcodes"`
}
