package dto

type PaymentMethodResponse struct {
	Code         string `json:"code"`
	Journal      string `json:"journal"`
	CurrencyCode string `json:"currency_code"`
	Description  string `json:"description"`
}
