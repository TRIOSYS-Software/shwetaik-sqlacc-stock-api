package handlers

import (
	"net/url"
	"shwetaik-sqlacc-stock-api/internal/delivery/dto"
	"shwetaik-sqlacc-stock-api/internal/domain/entities"
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
			PriceTag:   *stockItemPrice.PriceTag,
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
		PriceTag:   *stockItemPrice.PriceTag,
		UOM:        stockItemPrice.UOM,
		Qty:        stockItemPrice.Qty,
		StockValue: stockItemPrice.StockValue,
	}
	return utils.SuccessResponse(c, "success", response)
}

func (s *StockItemPriceHandler) CreateStockItemPrice(c *fiber.Ctx) error {
	rawCode := c.Params("code")
	code, err := url.QueryUnescape(rawCode)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	var stockItemPriceDTO dto.StockItemPriceRequest
	if err := c.BodyParser(&stockItemPriceDTO); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	var stockItemPrice entities.STItemPrice = entities.STItemPrice{
		Code:       code,
		PriceTag:   &stockItemPriceDTO.PriceTag,
		UOM:        stockItemPriceDTO.UOM,
		Qty:        stockItemPriceDTO.Qty,
		StockValue: stockItemPriceDTO.StockValue,
		TagType:    "C",
	}
	if err := s.usecase.CreateStockItemPrice(code, &stockItemPrice); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	response := dto.StockItemPriceResponse{
		DtlKey:     stockItemPrice.DtlKey,
		Code:       stockItemPrice.Code,
		PriceTag:   *stockItemPrice.PriceTag,
		UOM:        stockItemPrice.UOM,
		Qty:        stockItemPrice.Qty,
		StockValue: stockItemPrice.StockValue,
	}
	return utils.SuccessResponse(c, "success", response)
}

func (s *StockItemPriceHandler) UpdateStockItemPrice(c *fiber.Ctx) error {
	rawCode := c.Params("code")
	code, err := url.QueryUnescape(rawCode)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	dtlKey, err := c.ParamsInt("dtlKey")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	var stockItemPriceDTO dto.StockItemPriceRequest
	if err := c.BodyParser(&stockItemPriceDTO); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	var stockItemPrice entities.STItemPrice = entities.STItemPrice{
		DtlKey:     dtlKey,
		Code:       code,
		PriceTag:   &stockItemPriceDTO.PriceTag,
		UOM:        stockItemPriceDTO.UOM,
		Qty:        stockItemPriceDTO.Qty,
		StockValue: stockItemPriceDTO.StockValue,
		TagType:    "C",
	}
	if err := s.usecase.UpdateStockItemPrice(code, &stockItemPrice); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	response := dto.StockItemPriceResponse{
		DtlKey:     stockItemPrice.DtlKey,
		Code:       stockItemPrice.Code,
		PriceTag:   *stockItemPrice.PriceTag,
		UOM:        stockItemPrice.UOM,
		Qty:        stockItemPrice.Qty,
		StockValue: stockItemPrice.StockValue,
	}
	return utils.SuccessResponse(c, "success", response)
}
