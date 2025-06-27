package usecases

import (
	"shwetaik-sqlacc-stock-api/internal/domain/entities"
	"shwetaik-sqlacc-stock-api/internal/domain/repositories"
)

type StockItemUseCase struct {
	repo repositories.StockItemRepository
}

func NewStockItemUseCase(repo repositories.StockItemRepository) *StockItemUseCase {
	return &StockItemUseCase{repo: repo}
}

func (u StockItemUseCase) GetAllStockItems(filter map[string]any) ([]entities.STItem, error) {
	return u.repo.GetAllStockItems(filter)
}

func (u StockItemUseCase) GetStockItemByCode(code string) (*entities.STItem, error) {
	return u.repo.GetStockItemByCode(code)
}
