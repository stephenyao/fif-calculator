package holdingshandler

import (
	"fif-calculator/internal/model"
	"fif-calculator/internal/viewmodel"
	"fif-calculator/views/holdings"
	"net/http"
)

func (h *HoldingsHandler) List(w http.ResponseWriter, r *http.Request) {
	allHoldings, err := h.Repo.AllHoldings()

	if err != nil {
		http.Error(w, "Failed to get trades", http.StatusInternalServerError)
	}

	viewModels := convertToHoldingViewModels(allHoldings)

	err = holdings.HoldingsList(r.URL.Path, viewModels).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Could not render start page", http.StatusInternalServerError)
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
