package repositories

import "shwetaik-sqlacc-stock-api/internal/domain/entities"

type PaymentRepository interface {
	Create(payment *entities.Payment) error
}
