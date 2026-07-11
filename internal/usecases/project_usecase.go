package usecases

import (
	"shwetaik-sqlacc-stock-api/internal/delivery/dto"
	"shwetaik-sqlacc-stock-api/internal/domain/entities"
	"shwetaik-sqlacc-stock-api/internal/domain/repositories"
)

type ProjectUseCase struct {
	repo repositories.ProjectRepository
}

func NewProjectUseCase(repo repositories.ProjectRepository) *ProjectUseCase {
	return &ProjectUseCase{repo: repo}
}

func toProjectResponse(project entities.Project) *dto.ProjectResponse {
	response := &dto.ProjectResponse{
		Code:         project.Code,
		ProjectValue: project.ProjectValue,
		ProjectCost:  project.ProjectCost,
		IsActive:     project.IsActive,
	}
	if project.Description != nil {
		response.Description = *project.Description
	}
	if project.Description2 != nil {
		response.Description2 = *project.Description2
	}
	return response
}

func (u ProjectUseCase) GetAllProjects(filter map[string]any) ([]*dto.ProjectResponse, error) {
	projects, err := u.repo.GetAllProjects(filter)
	if err != nil {
		return nil, err
	}

	response := make([]*dto.ProjectResponse, 0, len(projects))
	for _, project := range projects {
		response = append(response, toProjectResponse(project))
	}
	return response, nil
}

func (u ProjectUseCase) GetProjectByCode(code string) (*dto.ProjectResponse, error) {
	project, err := u.repo.GetProjectByCode(code)
	if err != nil {
		return nil, err
	}
	return toProjectResponse(*project), nil
}
