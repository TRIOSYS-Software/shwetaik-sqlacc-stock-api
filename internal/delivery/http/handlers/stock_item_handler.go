package handlers

import (
	"net/url"
	"shwetaik-sqlacc-stock-api/internal/usecases"
	"shwetaik-sqlacc-stock-api/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

const defaultStockItemsLimit = 100
const maxStockItemsLimit = 1000

type StockItemHandler struct {
	usecase *usecases.StockItemUseCase
}

func NewStockItemHandler(usecase *usecases.StockItemUseCase) *StockItemHandler {
	return &StockItemHandler{usecase: usecase}
}

func (s *StockItemHandler) GetAllStockItems(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", defaultStockItemsLimit)
	if limit > maxStockItemsLimit {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "limit exceeds maximum allowed value")
	}

	filter := make(map[string]any)
	filter["limit"] = limit
	filter["after"] = c.Query("after", "")
	filter["stock_group"] = c.Query("stock_group", "")
	filter["description"] = c.Query("description", "")
	filter["location"] = c.Query("location", "")

	response, err := s.usecase.GetAllStockItems(filter)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	pagination := utils.Pagination{
		Limit:   limit,
		HasMore: len(response) == limit,
	}
	if len(response) > 0 {
		pagination.After = response[len(response)-1].Code
	}

	return utils.SuccessPaginatedResponse(c, "success", response, pagination)
}

func (s *StockItemHandler) GetStockItemByCode(c *fiber.Ctx) error {
	rawCode := c.Params("code")
	code, err := url.QueryUnescape(rawCode)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	location := c.Query("location", "")
	response, err := s.usecase.GetStockItemByCode(code, location)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return utils.SuccessResponse(c, "success", response)
}
