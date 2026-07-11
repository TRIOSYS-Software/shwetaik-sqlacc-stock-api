package handlers

import (
	"net/url"
	"shwetaik-sqlacc-stock-api/internal/usecases"
	"shwetaik-sqlacc-stock-api/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

const defaultGLAccLimit = 100
const maxGLAccLimit = 1000

type GLAccHandler struct {
	usecase *usecases.GLAccUseCase
}

func NewGLAccHandler(usecase *usecases.GLAccUseCase) *GLAccHandler {
	return &GLAccHandler{usecase: usecase}
}

func (h *GLAccHandler) GetAllGLAccs(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", defaultGLAccLimit)
	if limit > maxGLAccLimit {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "limit exceeds maximum allowed value")
	}

	filter := make(map[string]any)
	filter["limit"] = limit
	filter["after"] = c.Query("after", "")
	filter["acctype"] = c.Query("acctype", "")
	filter["description"] = c.Query("description", "")
	if parent := c.QueryInt("parent", 0); parent > 0 {
		filter["parent"] = parent
	}

	response, err := h.usecase.GetAllGLAccs(filter)
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

func (h *GLAccHandler) GetGLAccByCode(c *fiber.Ctx) error {
	rawCode := c.Params("code")
	code, err := url.QueryUnescape(rawCode)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	response, err := h.usecase.GetGLAccByCode(code)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return utils.SuccessResponse(c, "success", response)
}
