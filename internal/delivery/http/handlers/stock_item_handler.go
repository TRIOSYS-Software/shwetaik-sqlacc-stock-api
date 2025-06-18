package handlers

import (
	"net/url"
	"shwetaik-sqlacc-stock-api/internal/delivery/dto"
	"shwetaik-sqlacc-stock-api/internal/usecases"
	"shwetaik-sqlacc-stock-api/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type StockItemHandler struct {
	usecase *usecases.StockItemUseCase
}

func NewStockItemHandler(usecase *usecases.StockItemUseCase) *StockItemHandler {
	return &StockItemHandler{usecase: usecase}
}

func (s *StockItemHandler) GetAllStockItems(c *fiber.Ctx) error {
	limit := c.QueryInt("limit")
	offset := c.QueryInt("offset")
	stockItems, err := s.usecase.GetAllStockItems(limit, offset)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	var response []dto.StockItemResponse

	for _, stockItem := range stockItems {
		var stockItemPricesDTO []dto.StockItemPriceResponse
		for _, stockItemPrice := range stockItem.STItemPrices {
			stockItemPricesDTO = append(stockItemPricesDTO, dto.StockItemPriceResponse{
				DtlKey:     stockItemPrice.DtlKey,
				Code:       stockItemPrice.Code,
				PriceTag:   *stockItemPrice.PriceTag,
				UOM:        stockItemPrice.UOM,
				Qty:        stockItemPrice.Qty,
				StockValue: stockItemPrice.StockValue,
			})
		}
		response = append(response, dto.StockItemResponse{
			DocKey:       stockItem.DocKey,
			Code:         stockItem.Code,
			Description:  *stockItem.Description,
			StockGroup:   stockItem.StockGroup,
			STItemPrices: stockItemPricesDTO,
		})
	}
	return utils.SuccessResponse(c, "success", response)
}

func (s *StockItemHandler) GetStockItemByCode(c *fiber.Ctx) error {
	rawCode := c.Params("code")
	code, err := url.QueryUnescape(rawCode)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	stockItem, err := s.usecase.GetStockItemByCode(code)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	var stockItemPriceDTO []dto.StockItemPriceResponse
	for _, stockItemPrice := range stockItem.STItemPrices {
		stockItemPriceDTO = append(stockItemPriceDTO, dto.StockItemPriceResponse{
			DtlKey:     stockItemPrice.DtlKey,
			Code:       stockItemPrice.Code,
			PriceTag:   *stockItemPrice.PriceTag,
			UOM:        stockItemPrice.UOM,
			Qty:        stockItemPrice.Qty,
			StockValue: stockItemPrice.StockValue,
		})
	}
	response := dto.StockItemResponse{
		DocKey:       stockItem.DocKey,
		Code:         stockItem.Code,
		Description:  *stockItem.Description,
		StockGroup:   stockItem.StockGroup,
		STItemPrices: stockItemPriceDTO,
	}
	return utils.SuccessResponse(c, "success", response)
}
