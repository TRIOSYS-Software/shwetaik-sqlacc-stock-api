package repositories

// StockItemPriceGateway wraps the vendor SQL Account REST API's stock item
// endpoint for writing customer prices. The vendor has no standalone price
// endpoint — prices are a nested array (sdscustomerprice) inside the stock
// item resource, addressed by the item's DOCKEY.
type StockItemPriceGateway interface {
	PutStockItemPrice(dockey int, payload map[string]any) (map[string]any, error)
}
