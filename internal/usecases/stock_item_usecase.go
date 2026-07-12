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
		item := &dto.StockItemResponse{
			DocKey:         stockItem.DocKey,
			Code:           stockItem.Code,
			Description:    *stockItem.Description,
			StockGroup:     stockItem.StockGroup,
			Balance:        stockItem.Balance,
			STItemBarcodes: stockItemBarcodesDTO,
			STItemPrices:   stockItemPricesDTO,
		}
		if stockItem.Location != nil {
			item.Location = *stockItem.Location
		}
		response = append(response, item)
	}
	return response, nil
}

func (u StockItemUseCase) GetStockItemByCode(code string, location string) (*dto.StockItemResponse, error) {
	stockItem, err := u.repo.GetStockItemByCode(code, location)
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
		Balance:        stockItem.Balance,
		STItemPrices:   stockItemPriceDTO,
		STItemBarcodes: stockItemBarcodeDTO,
	}
	if stockItem.Location != nil {
		response.Location = *stockItem.Location
	}
	return &response, nil
}
