package holdingshandler

import (
	"fif-calculator/internal/model"
	"fif-calculator/internal/repository"
	"fif-calculator/internal/viewmodel"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type HoldingsHandler struct {
	HoldingsRepository repository.HoldingsRepository
	TradeRepository    repository.TradeRepository
}

func NewHoldingsHandler(db *sqlx.DB) *HoldingsHandler {
	return &HoldingsHandler{
		HoldingsRepository: repository.NewHoldingsRepository(db),
		TradeRepository:    repository.NewTradeRepository(db),
	}
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

func (h *HoldingsHandler) Index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/holdings", http.StatusSeeOther)
}
