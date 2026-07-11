package dto

type PaymentVoucherDetail struct {
	Code        string  `json:"code"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Project     string  `json:"project"`
}

type PaymentVoucherRequest struct {
	DocKey        string                 `json:"dockey"`
	DocDate       string                 `json:"docdate"`
	PaymentMethod string                 `json:"paymentmethod"`
	Description   string                 `json:"description"`
	Project       string                 `json:"project"`
	DocAmt        float64                `json:"docamt"`
	Details       []PaymentVoucherDetail `json:"sdsdocdetail"`
}
