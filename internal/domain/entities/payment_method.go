package entities

type PaymentMethod struct {
	Code         string  `gorm:"column:CODE"`
	Journal      *string `gorm:"column:JOURNAL"`
	CurrencyCode *string `gorm:"column:CURRENCYCODE"`
	Description  *string `gorm:"column:DESCRIPTION"`
}
