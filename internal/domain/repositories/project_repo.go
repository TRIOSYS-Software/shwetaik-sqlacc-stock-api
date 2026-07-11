package repositories

import (
	"shwetaik-sqlacc-stock-api/internal/domain/entities"
)

type ProjectRepository interface {
	GetAllProjects(filter map[string]any) ([]entities.Project, error)
	GetProjectByCode(code string) (*entities.Project, error)
}
