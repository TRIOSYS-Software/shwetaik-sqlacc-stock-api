package usecases

import (
	"shwetaik-sqlacc-stock-api/internal/delivery/dto"
	"shwetaik-sqlacc-stock-api/internal/domain/repositories"
)

type PaymentMethodUseCase struct {
	repo repositories.PaymentMethodRepository
}

func NewPaymentMethodUseCase(repo repositories.PaymentMethodRepository) *PaymentMethodUseCase {
	return &PaymentMethodUseCase{repo: repo}
}

func (u PaymentMethodUseCase) GetAllPaymentMethods() ([]*dto.PaymentMethodResponse, error) {
	paymentMethods, err := u.repo.GetAllPaymentMethods()
	if err != nil {
		return nil, err
	}

	response := make([]*dto.PaymentMethodResponse, 0, len(paymentMethods))
	for _, paymentMethod := range paymentMethods {
		item := &dto.PaymentMethodResponse{
			Code: paymentMethod.Code,
		}
		if paymentMethod.Journal != nil {
			item.Journal = *paymentMethod.Journal
		}
		if paymentMethod.CurrencyCode != nil {
			item.CurrencyCode = *paymentMethod.CurrencyCode
		}
		if paymentMethod.Description != nil {
			item.Description = *paymentMethod.Description
		}
		response = append(response, item)
	}
	return response, nil
}
