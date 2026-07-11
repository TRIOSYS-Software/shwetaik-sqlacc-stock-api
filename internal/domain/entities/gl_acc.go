package entities

type GLAcc struct {
	DocKey         int     `json:"dockey" gorm:"column:DOCKEY;primaryKey;not null"`
	Parent         int     `json:"parent" gorm:"column:PARENT;not null"`
	Code           string  `json:"code" gorm:"column:CODE;size:10;not null"`
	Description    *string `json:"description,omitempty" gorm:"column:DESCRIPTION;size:200"`
	Description2   *string `json:"description2,omitempty" gorm:"column:DESCRIPTION2;size:200"`
	AccType        *string `json:"acc_type,omitempty" gorm:"column:ACCTYPE;size:2"`
	SpecialAccType *string `json:"special_acc_type,omitempty" gorm:"column:SPECIALACCTYPE;size:2"`
	Tax            *string `json:"tax,omitempty" gorm:"column:TAX;size:10"`
	CashFlowType   *int    `json:"cash_flow_type,omitempty" gorm:"column:CASHFLOWTYPE"`
	SIC            *string `json:"sic,omitempty" gorm:"column:SIC;size:10"`
}

func (GLAcc) TableName() string {
	return "GL_ACC"
}
