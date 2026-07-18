package entities

import "time"

type Payment struct {
	DocKey        int             `json:"dockey" gorm:"column:DOCKEY;primaryKey;not null"`
	DocNo         string          `json:"docno" gorm:"column:DOCNO"`
	DocType       string          `json:"doctype" gorm:"column:DOCTYPE"`
	DocDate       time.Time       `json:"docdate" gorm:"column:DOCDATE;autoCreateTime"`
	PostDate      time.Time       `json:"postdate" gorm:"column:POSTDATE;autoUpdateTime"`
	TaxDate       time.Time       `json:"taxdate" gorm:"column:TAXDATE;autoCreateTime"`
	Description   *string         `json:"description,omitempty" gorm:"column:DESCRIPTION"`
	PaymentMethod string          `json:"paymentmethod" gorm:"column:PAYMENTMETHOD;not null"`
	Area          string          `json:"area" gorm:"column:AREA;default:----"`
	Agent         string          `json:"agent" gorm:"column:AGENT;default:----"`
	Project       string          `json:"project" gorm:"column:PROJECT;default:----"`
	Journal       string          `json:"journal" gorm:"column:JOURNAL"`
	CurrencyCode  string          `json:"currencycode" gorm:"column:CURRENCYCODE;default:----"`
	CurrencyRate  float64         `json:"currencyrate" gorm:"column:CURRENCYRATE;default:1"`
	DocAmt        float64         `json:"docamt" gorm:"column:DOCAMT"`
	LocalDocAmt   float64         `json:"localdocamt" gorm:"column:LOCALDOCAMT"`
	GLTransID     int             `json:"-" gorm:"column:GLTRANSID"`
	Cancelled     bool            `json:"cancelled" gorm:"column:CANCELLED;default:0"`
	LastModified  uint            `json:"-" gorm:"column:LASTMODIFIED;default:0"`
	Details       []PaymentDetail `json:"details,omitempty" gorm:"foreignKey:DocKey;references:DocKey"`
}

func (Payment) TableName() string {
	return "GL_CB"
}
