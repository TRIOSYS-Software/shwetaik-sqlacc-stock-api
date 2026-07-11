package repositories

import (
	"shwetaik-sqlacc-stock-api/internal/domain/entities"
)

type PaymentMethodRepository interface {
	GetAllPaymentMethods() ([]entities.PaymentMethod, error)
}
