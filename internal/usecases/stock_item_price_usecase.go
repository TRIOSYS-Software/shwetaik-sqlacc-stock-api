package usecases

import (
	"shwetaik-sqlacc-stock-api/internal/delivery/dto"
	"shwetaik-sqlacc-stock-api/internal/domain/entities"
	"shwetaik-sqlacc-stock-api/internal/domain/repositories"
)

type StockItemPriceUseCase struct {
	repo repositories.StockItemPriceRepository
}

func NewStockItemPriceUseCase(repo repositories.StockItemPriceRepository) *StockItemPriceUseCase {
	return &StockItemPriceUseCase{repo: repo}
}

func (u StockItemPriceUseCase) GetStockItemPricesByCode(code string) ([]dto.StockItemPriceResponse, error) {
	stockItemPrices, err := u.repo.GetStockItemPricesByCode(code)
	if err != nil {
		return nil, err
	}
	var response []dto.StockItemPriceResponse

	for _, stockItemPrice := range stockItemPrices {
		response = append(response, dto.StockItemPriceResponse{
			DtlKey:     stockItemPrice.DtlKey,
			Code:       stockItemPrice.Code,
			PriceTag:   stockItemPrice.PriceTag,
			UOM:        stockItemPrice.UOM,
			Qty:        stockItemPrice.Qty,
			StockValue: stockItemPrice.StockValue,
		})
	}
	return response, nil
}

func (u StockItemPriceUseCase) GetStockItemPriceByDTLKey(code string, dtlKey int) (*dto.StockItemPriceResponse, error) {
	stockItemPrice, err := u.repo.GetStockItemPriceByDTLKey(code, dtlKey)
	if err != nil {
		return nil, err
	}
	response := dto.StockItemPriceResponse{
		DtlKey:     stockItemPrice.DtlKey,
		Code:       stockItemPrice.Code,
		PriceTag:   stockItemPrice.PriceTag,
		UOM:        stockItemPrice.UOM,
		Qty:        stockItemPrice.Qty,
		StockValue: stockItemPrice.StockValue,
	}
	return &response, nil
}

func (u StockItemPriceUseCase) CreateStockItemPrice(code string, stockItemPriceDTO dto.StockItemPriceRequest) (*dto.StockItemPriceResponse, error) {
	var stockItemPrice entities.STItemPrice = entities.STItemPrice{
		Code:       code,
		PriceTag:   &stockItemPriceDTO.PriceTag,
		UOM:        stockItemPriceDTO.UOM,
		Qty:        stockItemPriceDTO.Qty,
		StockValue: stockItemPriceDTO.StockValue,
		TagType:    "C",
	}
	if err := u.repo.CreateStockItemPrice(&stockItemPrice); err != nil {
		return nil, err
	}
	response := dto.StockItemPriceResponse{
		DtlKey:     stockItemPrice.DtlKey,
		Code:       stockItemPrice.Code,
		PriceTag:   stockItemPrice.PriceTag,
		UOM:        stockItemPrice.UOM,
		Qty:        stockItemPrice.Qty,
		StockValue: stockItemPrice.StockValue,
	}
	return &response, nil
}

func (u StockItemPriceUseCase) UpdateStockItemPrice(code string, dtlKey int, stockItemPriceDTO *dto.StockItemPriceRequest) (*dto.StockItemPriceResponse, error) {
	var stockItemPrice entities.STItemPrice = entities.STItemPrice{
		DtlKey:     dtlKey,
		Code:       code,
		PriceTag:   &stockItemPriceDTO.PriceTag,
		UOM:        stockItemPriceDTO.UOM,
		Qty:        stockItemPriceDTO.Qty,
		StockValue: stockItemPriceDTO.StockValue,
		TagType:    "C",
	}
	if err := u.repo.UpdateStockItemPrice(code, &stockItemPrice); err != nil {
		return nil, err
	}
	response := dto.StockItemPriceResponse{
		DtlKey:     stockItemPrice.DtlKey,
		Code:       stockItemPrice.Code,
		PriceTag:   stockItemPrice.PriceTag,
		UOM:        stockItemPrice.UOM,
		Qty:        stockItemPrice.Qty,
		StockValue: stockItemPrice.StockValue,
	}
	return &response, nil
}
