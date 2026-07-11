package handlers

import (
	"shwetaik-sqlacc-stock-api/internal/delivery/dto"
	"shwetaik-sqlacc-stock-api/internal/usecases"
	"shwetaik-sqlacc-stock-api/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type PaymentVoucherHandler struct {
	usecase *usecases.PaymentVoucherUseCase
}

func NewPaymentVoucherHandler(usecase *usecases.PaymentVoucherUseCase) *PaymentVoucherHandler {
	return &PaymentVoucherHandler{usecase: usecase}
}

func (h *PaymentVoucherHandler) CreatePaymentVoucher(c *fiber.Ctx) error {
	var req dto.PaymentVoucherRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	response, err := h.usecase.CreateExpensePaymentVoucher(req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadGateway, err.Error())
	}
	return utils.SuccessResponse(c, "success", response)
}
