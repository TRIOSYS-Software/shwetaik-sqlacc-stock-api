package handlers

import (
	"net/url"
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
	filter := make(map[string]any)
	filter["limit"] = c.QueryInt("limit")
	filter["offset"] = c.QueryInt("offset")
	filter["stock_group"] = c.Query("stock_group", "")
	filter["description"] = c.Query("description", "")

	response, err := s.usecase.GetAllStockItems(filter)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return utils.SuccessResponse(c, "success", response)
}

func (s *StockItemHandler) GetStockItemByCode(c *fiber.Ctx) error {
	rawCode := c.Params("code")
	code, err := url.QueryUnescape(rawCode)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	response, err := s.usecase.GetStockItemByCode(code)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return utils.SuccessResponse(c, "success", response)
}
