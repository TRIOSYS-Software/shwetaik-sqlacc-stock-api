package usecases

import (
	"fmt"
	"time"

	"shwetaik-sqlacc-stock-api/internal/delivery/dto"
	"shwetaik-sqlacc-stock-api/internal/domain/entities"
	"shwetaik-sqlacc-stock-api/internal/domain/repositories"
)

type PaymentUseCase struct {
	repo repositories.PaymentRepository
}

func NewPaymentUseCase(repo repositories.PaymentRepository) *PaymentUseCase {
	return &PaymentUseCase{repo: repo}
}

// docDateLayouts are the formats accepted for PaymentVoucherRequest.DocDate,
// tried in order.
var docDateLayouts = []string{time.RFC3339, "2006-01-02"}

func parseDocDate(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	var lastErr error
	for _, layout := range docDateLayouts {
		if t, err := time.Parse(layout, value); err == nil {
			return t, nil
		} else {
			lastErr = err
		}
	}
	return time.Time{}, fmt.Errorf("docdate %q is not a valid date (expected RFC3339 or YYYY-MM-DD): %w", value, lastErr)
}

// CreatePayment accepts the same request shape as the vendor-API-backed
// POST /payment-vouchers (dto.PaymentVoucherRequest) but writes directly
// into GL_CB/GL_CBDTL/GL_TRANS instead. DocAmt is ignored — it's always
// recomputed server-side as the sum of the detail amounts, since the ledger
// must balance regardless of what the client sends.
func (u PaymentUseCase) CreatePayment(req dto.PaymentVoucherRequest) (*dto.PaymentResponse, error) {
	docDate, err := parseDocDate(req.DocDate)
	if err != nil {
		return nil, err
	}

	details := make([]entities.PaymentDetail, 0, len(req.Details))
	for _, d := range req.Details {
		description := d.Description
		details = append(details, entities.PaymentDetail{
			Code:        d.Code,
			Description: &description,
			Amount:      d.Amount,
			// No FX conversion in this request shape — local and currency
			// amounts equal the transaction amount (CurrencyRate 1:1).
			LocalAmount:    d.Amount,
			CurrencyAmount: d.Amount,
			Project:        d.Project,
		})
	}

	description := req.Description
	payment := &entities.Payment{
		DocNo:         req.DocNo,
		DocDate:       docDate,
		PaymentMethod: req.PaymentMethod,
		Description:   &description,
		Project:       req.Project,
		Details:       details,
	}

	if err := u.repo.Create(payment); err != nil {
		return nil, err
	}

	// payment.Details is nilled out by repo.Create before the final save, so
	// the response is built from the local `details` slice — repo.Create
	// mutated its elements (DtlKey, etc.) in place via the same backing array.
	responseDetails := make([]dto.PaymentDetailResponse, 0, len(details))
	for _, d := range details {
		desc := ""
		if d.Description != nil {
			desc = *d.Description
		}
		responseDetails = append(responseDetails, dto.PaymentDetailResponse{
			DtlKey:      d.DtlKey,
			Code:        d.Code,
			Description: desc,
			Amount:      d.Amount,
		})
	}

	return &dto.PaymentResponse{
		DocKey:        payment.DocKey,
		DocNo:         payment.DocNo,
		PaymentMethod: payment.PaymentMethod,
		Description:   req.Description,
		Project:       payment.Project,
		DocAmt:        payment.DocAmt,
		Details:       responseDetails,
	}, nil
}
