package dto

type StockItemResponse struct {
	DocKey       int                      `json:"dockey"`
	Code         string                   `json:"code"`
	Description  string                   `json:"description"`
	StockGroup   string                   `json:"stock_group"`
	STItemPrices []StockItemPriceResponse `json:"stock_item_prices"`
}
