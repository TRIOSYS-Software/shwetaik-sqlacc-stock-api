package usecases

import (
	"shwetaik-sqlacc-stock-api/internal/delivery/dto"
	"shwetaik-sqlacc-stock-api/internal/domain/repositories"
)

type StockItemUseCase struct {
	repo repositories.StockItemRepository
}

func NewStockItemUseCase(repo repositories.StockItemRepository) *StockItemUseCase {
	return &StockItemUseCase{repo: repo}
}

func (u StockItemUseCase) GetAllStockItems(filter map[string]any) ([]*dto.StockItemResponse, error) {
	stockItems, err := u.repo.GetAllStockItems(filter)
	if err != nil {
		return nil, err
	}

	var response []*dto.StockItemResponse

	for _, stockItem := range stockItems {
		var stockItemPricesDTO []dto.StockItemPriceResponse
		for _, stockItemPrice := range stockItem.STItemPrices {
			stockItemPricesDTO = append(stockItemPricesDTO, dto.StockItemPriceResponse{
				DtlKey:     stockItemPrice.DtlKey,
				Code:       stockItemPrice.Code,
				PriceTag:   stockItemPrice.PriceTag,
				UOM:        stockItemPrice.UOM,
				Qty:        stockItemPrice.Qty,
				StockValue: stockItemPrice.StockValue,
			})
		}
		var stockItemBarcodesDTO []dto.StockItemBarcodeResponse
		for _, stockItemBarcode := range stockItem.STItemBarcodes {
			stockItemBarcodesDTO = append(stockItemBarcodesDTO, dto.StockItemBarcodeResponse{
				AutoKey: stockItemBarcode.AutoKey,
				Barcode: stockItemBarcode.Barcode,
				UOM:     stockItemBarcode.UOM,
			})
		}
		response = append(response, &dto.StockItemResponse{
			DocKey:         stockItem.DocKey,
			Code:           stockItem.Code,
			Description:    *stockItem.Description,
			StockGroup:     stockItem.StockGroup,
			STItemBarcodes: stockItemBarcodesDTO,
			STItemPrices:   stockItemPricesDTO,
		})
	}
	return response, nil
}

func (u StockItemUseCase) GetStockItemByCode(code string) (*dto.StockItemResponse, error) {
	stockItem, err := u.repo.GetStockItemByCode(code)
	if err != nil {
		return nil, err
	}
	var stockItemPriceDTO []dto.StockItemPriceResponse
	for _, stockItemPrice := range stockItem.STItemPrices {
		stockItemPriceDTO = append(stockItemPriceDTO, dto.StockItemPriceResponse{
			DtlKey:     stockItemPrice.DtlKey,
			Code:       stockItemPrice.Code,
			PriceTag:   stockItemPrice.PriceTag,
			UOM:        stockItemPrice.UOM,
			Qty:        stockItemPrice.Qty,
			StockValue: stockItemPrice.StockValue,
		})
	}
	var stockItemBarcodeDTO []dto.StockItemBarcodeResponse
	for _, stockItemBarcode := range stockItem.STItemBarcodes {
		stockItemBarcodeDTO = append(stockItemBarcodeDTO, dto.StockItemBarcodeResponse{
			AutoKey: stockItemBarcode.AutoKey,
			Barcode: stockItemBarcode.Barcode,
			UOM:     stockItemBarcode.UOM,
		})
	}
	response := dto.StockItemResponse{
		DocKey:         stockItem.DocKey,
		Code:           stockItem.Code,
		Description:    *stockItem.Description,
		StockGroup:     stockItem.StockGroup,
		STItemPrices:   stockItemPriceDTO,
		STItemBarcodes: stockItemBarcodeDTO,
	}
	return &response, nil
}
