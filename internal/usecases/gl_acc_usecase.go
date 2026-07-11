package usecases

import (
	"shwetaik-sqlacc-stock-api/internal/delivery/dto"
	"shwetaik-sqlacc-stock-api/internal/domain/entities"
	"shwetaik-sqlacc-stock-api/internal/domain/repositories"
)

type GLAccUseCase struct {
	repo repositories.GLAccRepository
}

func NewGLAccUseCase(repo repositories.GLAccRepository) *GLAccUseCase {
	return &GLAccUseCase{repo: repo}
}

func toGLAccResponse(glAcc entities.GLAcc) *dto.GLAccResponse {
	response := &dto.GLAccResponse{
		DocKey:       glAcc.DocKey,
		Parent:       glAcc.Parent,
		Code:         glAcc.Code,
		CashFlowType: glAcc.CashFlowType,
	}
	if glAcc.Description != nil {
		response.Description = *glAcc.Description
	}
	if glAcc.Description2 != nil {
		response.Description2 = *glAcc.Description2
	}
	if glAcc.AccType != nil {
		response.AccType = *glAcc.AccType
	}
	if glAcc.SpecialAccType != nil {
		response.SpecialAccType = *glAcc.SpecialAccType
	}
	if glAcc.Tax != nil {
		response.Tax = *glAcc.Tax
	}
	if glAcc.SIC != nil {
		response.SIC = *glAcc.SIC
	}
	return response
}

func (u GLAccUseCase) GetAllGLAccs(filter map[string]any) ([]*dto.GLAccResponse, error) {
	glAccs, err := u.repo.GetAllGLAccs(filter)
	if err != nil {
		return nil, err
	}

	response := make([]*dto.GLAccResponse, 0, len(glAccs))
	for _, glAcc := range glAccs {
		response = append(response, toGLAccResponse(glAcc))
	}
	return response, nil
}

func (u GLAccUseCase) GetGLAccByCode(code string) (*dto.GLAccResponse, error) {
	glAcc, err := u.repo.GetGLAccByCode(code)
	if err != nil {
		return nil, err
	}
	return toGLAccResponse(*glAcc), nil
}
