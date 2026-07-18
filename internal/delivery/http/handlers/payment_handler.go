package handlers

import (
	"shwetaik-sqlacc-stock-api/internal/delivery/dto"
	"shwetaik-sqlacc-stock-api/internal/usecases"
	"shwetaik-sqlacc-stock-api/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type PaymentHandler struct {
	usecase *usecases.PaymentUseCase
}

func NewPaymentHandler(usecase *usecases.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{usecase: usecase}
}

func (h *PaymentHandler) CreatePayment(c *fiber.Ctx) error {
	var req dto.PaymentVoucherRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	response, err := h.usecase.CreatePayment(req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return utils.SuccessResponse(c, "success", response)
}
