package usecases

import (
	"context"

	"shwetaik-sqlacc-stock-api/internal/delivery/dto"
	"shwetaik-sqlacc-stock-api/internal/domain/entities"
	"shwetaik-sqlacc-stock-api/internal/domain/repositories"
	"shwetaik-sqlacc-stock-api/pkg/utils"
)

type StockItemPriceUseCase struct {
	repo          repositories.StockItemPriceRepository
	stockItemRepo repositories.StockItemRepository
	vendorGateway repositories.StockItemPriceGateway
}

func NewStockItemPriceUseCase(
	repo repositories.StockItemPriceRepository,
	stockItemRepo repositories.StockItemRepository,
	vendorGateway repositories.StockItemPriceGateway,
) *StockItemPriceUseCase {
	return &StockItemPriceUseCase{repo: repo, stockItemRepo: stockItemRepo, vendorGateway: vendorGateway}
}

func (u StockItemPriceUseCase) GetStockItemPricesByCode(code string) ([]entities.STItemPrice, error) {
	return u.repo.GetStockItemPricesByCode(code)
}

func (u StockItemPriceUseCase) GetStockItemPriceByDTLKey(code string, dtlKey int) (*entities.STItemPrice, error) {
	return u.repo.GetStockItemPriceByDTLKey(code, dtlKey)
}

// PutStockItemPrices writes through the vendor SQL Account API instead of
// the database — the vendor has no standalone price endpoint, so this puts
// the whole customer price list onto the stock item's sdscustomerprice
// array via the vendor's stock item PUT, addressed by the item's DOCKEY
// (looked up from our own DB). This is a true replace: the vendor swaps out
// the entire array on each PUT, so the caller must include every price
// line it wants kept. A line with no dtlkey is created as new; a line with
// a dtlkey replaces the existing line with that key.
//
// The vendor responds with the entire stock item record (BOM, barcodes,
// base64 picture, everything), so only sdscustomerprice is picked back out
// of that response rather than relaying it whole.
func (u StockItemPriceUseCase) PutStockItemPrices(ctx context.Context, code string, items []dto.StockItemPriceItem) ([]dto.StockItemPriceResponse, error) {
	stockItem, err := u.stockItemRepo.GetStockItemByCode(code)
	if err != nil {
		return nil, err
	}

	priceLines := make([]map[string]any, 0, len(items))
	for _, item := range items {
		line := map[string]any{
			"pricetag":   item.PriceTag,
			"uom":        item.UOM,
			"qty":        item.Qty,
			"stockvalue": item.StockValue,
		}
		if item.DtlKey != nil {
			line["dtlkey"] = *item.DtlKey
		}
		priceLines = append(priceLines, line)
	}

	payload := map[string]any{
		"code":             code,
		"sdscustomerprice": priceLines,
	}

	vendorResponse, err := u.vendorGateway.PutStockItemPrice(ctx, stockItem.DocKey, payload)
	if err != nil {
		return nil, err
	}

	return extractCustomerPrices(code, vendorResponse), nil
}

func extractCustomerPrices(code string, vendorResponse map[string]any) []dto.StockItemPriceResponse {
	response := []dto.StockItemPriceResponse{}

	rawPrices, ok := vendorResponse["sdscustomerprice"].([]any)
	if !ok {
		return response
	}

	for _, rawPrice := range rawPrices {
		priceMap, ok := rawPrice.(map[string]any)
		if !ok {
			continue
		}
		response = append(response, dto.StockItemPriceResponse{
			DtlKey:     utils.AnyToInt(priceMap["dtlkey"]),
			Code:       code,
			PriceTag:   utils.AnyToStringPtr(priceMap["pricetag"]),
			UOM:        utils.AnyToString(priceMap["uom"]),
			Qty:        utils.AnyToFloatPtr(priceMap["qty"]),
			StockValue: utils.AnyToFloatPtr(priceMap["stockvalue"]),
		})
	}
	return response
}
