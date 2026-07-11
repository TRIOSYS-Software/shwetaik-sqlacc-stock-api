package dto

type ProjectResponse struct {
	Code         string   `json:"code"`
	Description  string   `json:"description,omitempty"`
	Description2 string   `json:"description2,omitempty"`
	ProjectValue *float64 `json:"project_value,omitempty"`
	ProjectCost  *float64 `json:"project_cost,omitempty"`
	IsActive     *bool    `json:"is_active,omitempty"`
}
