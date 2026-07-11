package handlers

import (
	"net/url"
	"shwetaik-sqlacc-stock-api/internal/usecases"
	"shwetaik-sqlacc-stock-api/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

const defaultProjectsLimit = 100
const maxProjectsLimit = 1000

type ProjectHandler struct {
	usecase *usecases.ProjectUseCase
}

func NewProjectHandler(usecase *usecases.ProjectUseCase) *ProjectHandler {
	return &ProjectHandler{usecase: usecase}
}

func (h *ProjectHandler) GetAllProjects(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", defaultProjectsLimit)
	if limit > maxProjectsLimit {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "limit exceeds maximum allowed value")
	}

	filter := make(map[string]any)
	filter["limit"] = limit
	filter["after"] = c.Query("after", "")

	response, err := h.usecase.GetAllProjects(filter)
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

func (h *ProjectHandler) GetProjectByCode(c *fiber.Ctx) error {
	rawCode := c.Params("code")
	code, err := url.QueryUnescape(rawCode)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	response, err := h.usecase.GetProjectByCode(code)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return utils.SuccessResponse(c, "success", response)
}
