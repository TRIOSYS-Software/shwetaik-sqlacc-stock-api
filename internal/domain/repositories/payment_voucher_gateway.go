package repositories

import "context"

// PaymentVoucherGateway wraps the vendor SQL Account REST API's payment
// voucher endpoint for write/create operations. Unlike the other
// repositories in this package, it does not talk to the Firebird database
// directly — it proxies to the vendor's own API.
type PaymentVoucherGateway interface {
	CreatePaymentVoucher(ctx context.Context, payload map[string]any) (map[string]any, error)
}
