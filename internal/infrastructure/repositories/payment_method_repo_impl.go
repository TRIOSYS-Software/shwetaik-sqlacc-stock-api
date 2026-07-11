package repositories

import (
	"shwetaik-sqlacc-stock-api/internal/domain/entities"
	"shwetaik-sqlacc-stock-api/internal/domain/repositories"

	"gorm.io/gorm"
)

type PaymentMethodRepositoryImpl struct {
	db *gorm.DB
}

func NewPaymentMethodRepository(db *gorm.DB) repositories.PaymentMethodRepository {
	return &PaymentMethodRepositoryImpl{db: db}
}

func (r *PaymentMethodRepositoryImpl) GetAllPaymentMethods() ([]entities.PaymentMethod, error) {
	// Only the columns the response needs are selected — PMMETHOD has BLOB
	// columns (GIRO, "DATA", ATTACHMENTS) that aren't used here.
	query := `
		SELECT pm.CODE, pm.JOURNAL, pm.CURRENCYCODE, gl.DESCRIPTION
		FROM PMMETHOD pm
		JOIN GL_ACC gl ON pm.CODE = gl.CODE
		ORDER BY pm.CODE
	`
	var paymentMethods []entities.PaymentMethod
	if err := r.db.Raw(query).Scan(&paymentMethods).Error; err != nil {
		return nil, err
	}
	return paymentMethods, nil
}
