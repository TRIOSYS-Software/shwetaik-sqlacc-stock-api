package repositories

import (
	"shwetaik-sqlacc-stock-api/internal/domain/entities"
)

type GLAccRepository interface {
	GetAllGLAccs(filter map[string]any) ([]entities.GLAcc, error)
	GetGLAccByCode(code string) (*entities.GLAcc, error)
}
