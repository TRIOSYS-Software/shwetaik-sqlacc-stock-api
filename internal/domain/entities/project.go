package entities

type Project struct {
	Code         string   `json:"code" gorm:"column:CODE;size:20;not null"`
	Description  *string  `json:"description,omitempty" gorm:"column:DESCRIPTION;size:80"`
	Description2 *string  `json:"description2,omitempty" gorm:"column:DESCRIPTION2;size:80"`
	ProjectValue *float64 `json:"project_value,omitempty" gorm:"column:PROJECTVALUE;precision:18;scale:2"`
	ProjectCost  *float64 `json:"project_cost,omitempty" gorm:"column:PROJECTCOST;precision:18;scale:2"`
	IsActive     *bool    `json:"is_active,omitempty" gorm:"column:ISACTIVE"`
}

func (Project) TableName() string {
	return "PROJECT"
}
