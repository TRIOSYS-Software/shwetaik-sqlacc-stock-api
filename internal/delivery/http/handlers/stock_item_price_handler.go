package handlers

import (
	"net/url"
	"shwetaik-sqlacc-stock-api/internal/delivery/dto"
	"shwetaik-sqlacc-stock-api/internal/usecases"
	"shwetaik-sqlacc-stock-api/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type StockItemPriceHandler struct {
	usecase *usecases.StockItemPriceUseCase
}

func NewStockItemPriceHandler(usecase *usecases.StockItemPriceUseCase) *StockItemPriceHandler {
	return &StockItemPriceHandler{usecase: usecase}
}

func (s *StockItemPriceHandler) GetStockItemPricesByCode(c *fiber.Ctx) error {
	rawCode := c.Params("code")
	code, err := url.QueryUnescape(rawCode)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	stockItemPrices, err := s.usecase.GetStockItemPricesByCode(code)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
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
	return utils.SuccessResponse(c, "success", response)
}

func (s *StockItemPriceHandler) GetStockItemPriceByDTLKey(c *fiber.Ctx) error {
	rawCode := c.Params("code")
	code, err := url.QueryUnescape(rawCode)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	dtlKey, err := c.ParamsInt("dtlKey")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	stockItemPrice, err := s.usecase.GetStockItemPriceByDTLKey(code, dtlKey)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	response := dto.StockItemPriceResponse{
		DtlKey:     stockItemPrice.DtlKey,
		Code:       stockItemPrice.Code,
		PriceTag:   stockItemPrice.PriceTag,
		UOM:        stockItemPrice.UOM,
		Qty:        stockItemPrice.Qty,
		StockValue: stockItemPrice.StockValue,
	}
	return utils.SuccessResponse(c, "success", response)
}

// PutStockItemPrices replaces the stock item's entire customer price list.
// This is a true replace, not a merge — the caller must include every
// price line it wants kept. A line with no dtlkey is created as new; a
// line with a dtlkey replaces the existing line with that key.
func (s *StockItemPriceHandler) PutStockItemPrices(c *fiber.Ctx) error {
	rawCode := c.Params("code")
	code, err := url.QueryUnescape(rawCode)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	var req dto.PutStockItemPricesRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	response, err := s.usecase.PutStockItemPrices(code, req.Prices)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return utils.SuccessResponse(c, "success", response)
}
