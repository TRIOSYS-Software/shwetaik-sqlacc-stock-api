package dto

type StockItemPriceResponse struct {
	DtlKey   int      `json:"dtl_key"`
	Code     string   `json:"code"`
	PriceTag string   `json:"price_tag"`
	UOM      string   `json:"uom"`
	Qty      *float64 `json:"qty"`
}
