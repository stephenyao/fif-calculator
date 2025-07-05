package holdingshandler

import (
	"fif-calculator/internal/model"
	"fif-calculator/internal/repository"
	"fif-calculator/internal/viewmodel"
	"github.com/jmoiron/sqlx"
)

type HoldingsHandler struct {
	Repo repository.HoldingsRepository
}

func NewHoldingsHandler(db *sqlx.DB) *HoldingsHandler {
	return &HoldingsHandler{Repo: repository.NewHoldingsRepository(db)}
}

func convertToHoldingViewModels(records []*model.HoldingRecord) []viewmodel.HoldingViewModel {
	var viewModels []viewmodel.HoldingViewModel
	for _, r := range records {
		viewModels = append(viewModels, viewmodel.HoldingViewModel{
			ID:       r.ID,
			Name:     r.Name,
			Symbol:   r.Symbol,
			Currency: r.Currency,
		})
	}
	return viewModels
}
