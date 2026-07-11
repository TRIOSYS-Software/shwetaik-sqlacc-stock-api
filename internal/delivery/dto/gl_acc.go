package dto

type GLAccResponse struct {
	DocKey         int     `json:"dockey"`
	Parent         int     `json:"parent"`
	Code           string  `json:"code"`
	Description    string  `json:"description,omitempty"`
	Description2   string  `json:"description2,omitempty"`
	AccType        string  `json:"acc_type,omitempty"`
	SpecialAccType string  `json:"special_acc_type,omitempty"`
	Tax            string  `json:"tax,omitempty"`
	CashFlowType   *int    `json:"cash_flow_type,omitempty"`
	SIC            string  `json:"sic,omitempty"`
}
