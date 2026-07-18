package entities

type PaymentDetail struct {
	DtlKey         int     `json:"dtlkey" gorm:"column:DTLKEY;primaryKey;not null"`
	DocKey         int     `json:"dockey" gorm:"column:DOCKEY;not null"`
	Seq            uint    `json:"seq" gorm:"column:SEQ"`
	Area           string  `json:"area" gorm:"column:AREA;default:----"`
	Agent          string  `json:"agent" gorm:"column:AGENT;default:----"`
	Project        string  `json:"project" gorm:"column:PROJECT;default:----"`
	Code           string  `json:"code" gorm:"column:CODE;not null"`
	Description    *string `json:"description,omitempty" gorm:"column:DESCRIPTION"`
	TaxInclusive   bool    `json:"tax_inclusive" gorm:"column:TAXINCLUSIVE"`
	Amount         float64 `json:"amount" gorm:"column:AMOUNT"`
	LocalAmount    float64 `json:"localamount" gorm:"column:LOCALAMOUNT"`
	CurrencyCode   string  `json:"currencycode" gorm:"column:CURRENCYCODE"`
	CurrencyRate   float64 `json:"currencyrate" gorm:"column:CURRENCYRATE"`
	CurrencyAmount float64 `json:"currencyamount" gorm:"column:CURRENCYAMOUNT"`
}

func (PaymentDetail) TableName() string {
	return "GL_CBDTL"
}
