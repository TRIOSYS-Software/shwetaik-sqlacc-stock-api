package usecases

import (
	"shwetaik-sqlacc-stock-api/internal/delivery/dto"
	"shwetaik-sqlacc-stock-api/internal/domain/repositories"
)

type PaymentVoucherUseCase struct {
	gateway repositories.PaymentVoucherGateway
}

func NewPaymentVoucherUseCase(gateway repositories.PaymentVoucherGateway) *PaymentVoucherUseCase {
	return &PaymentVoucherUseCase{gateway: gateway}
}

// CreateExpensePaymentVoucher wraps the vendor SQL Account API's payment
// voucher create endpoint. The field names below (docno, docdate,
// paymentmethod, docamt, sdsdocdetail, ...) match the vendor's actual
// schema, not the simplified request this endpoint accepts. docno is left
// as an empty string — the vendor assigns it.
func (u PaymentVoucherUseCase) CreateExpensePaymentVoucher(req dto.PaymentVoucherRequest) (map[string]any, error) {
	var sdsDetails []map[string]any
	for _, detail := range req.Details {
		sdsDetails = append(sdsDetails, map[string]any{
			"code":        detail.Code,
			"description": detail.Description,
			"amount":      detail.Amount,
			"project":     detail.Project,
		})
	}
	payload := map[string]any{
		"docno":         req.DocKey,
		"docdate":       req.DocDate,
		"paymentmethod": req.PaymentMethod,
		"description":   req.Description,
		"project":       req.Project,
		"docamt":        req.DocAmt,
		"sdsdocdetail":  sdsDetails,
	}

	return u.gateway.CreatePaymentVoucher(payload)
}
