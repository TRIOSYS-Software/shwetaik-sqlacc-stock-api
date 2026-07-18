package dto

type PaymentDetailResponse struct {
	DtlKey      int     `json:"dtlkey"`
	Code        string  `json:"code"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
}

type PaymentResponse struct {
	DocKey        int                     `json:"dockey"`
	DocNo         string                  `json:"docno"`
	PaymentMethod string                  `json:"paymentmethod"`
	Description   string                  `json:"description"`
	Project       string                  `json:"project"`
	DocAmt        float64                 `json:"docamt"`
	Details       []PaymentDetailResponse `json:"details"`
}
