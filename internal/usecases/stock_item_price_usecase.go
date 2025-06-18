package usecases

import (
	"shwetaik-sqlacc-stock-api/internal/domain/entities"
	"shwetaik-sqlacc-stock-api/internal/domain/repositories"
)

type StockItemPriceUseCase struct {
	repo repositories.StockItemPriceRepository
}

func NewStockItemPriceUseCase(repo repositories.StockItemPriceRepository) *StockItemPriceUseCase {
	return &StockItemPriceUseCase{repo: repo}
}

func (u StockItemPriceUseCase) GetStockItemPricesByCode(code string) ([]entities.STItemPrice, error) {
	return u.repo.GetStockItemPricesByCode(code)
}

func (u StockItemPriceUseCase) GetStockItemPriceByDTLKey(code string, dtlKey int) (*entities.STItemPrice, error) {
	return u.repo.GetStockItemPriceByDTLKey(code, dtlKey)
}

func (u StockItemPriceUseCase) CreateStockItemPrice(code string, stockItemPrice *entities.STItemPrice) error {
	return u.repo.CreateStockItemPrice(code, stockItemPrice)
}

func (u StockItemPriceUseCase) UpdateStockItemPrice(code string, stockItemPrice *entities.STItemPrice) error {
	return u.repo.UpdateStockItemPrice(code, stockItemPrice)
}
