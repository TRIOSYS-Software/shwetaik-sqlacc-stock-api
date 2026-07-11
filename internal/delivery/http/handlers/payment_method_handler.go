package handlers

import (
	"shwetaik-sqlacc-stock-api/internal/usecases"
	"shwetaik-sqlacc-stock-api/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type PaymentMethodHandler struct {
	usecase *usecases.PaymentMethodUseCase
}

func NewPaymentMethodHandler(usecase *usecases.PaymentMethodUseCase) *PaymentMethodHandler {
	return &PaymentMethodHandler{usecase: usecase}
}

func (h *PaymentMethodHandler) GetAllPaymentMethods(c *fiber.Ctx) error {
	response, err := h.usecase.GetAllPaymentMethods()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return utils.SuccessResponse(c, "success", response)
}
