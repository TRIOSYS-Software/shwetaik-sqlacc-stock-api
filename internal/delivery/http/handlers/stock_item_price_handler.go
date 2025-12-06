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
	response, err := s.usecase.GetStockItemPricesByCode(code)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
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
	response, err := s.usecase.GetStockItemPriceByDTLKey(code, dtlKey)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
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
	response, err := s.usecase.CreateStockItemPrice(code, stockItemPriceDTO)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
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
	response, err := s.usecase.UpdateStockItemPrice(code, dtlKey, &stockItemPriceDTO)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return utils.SuccessResponse(c, "success", response)
}
