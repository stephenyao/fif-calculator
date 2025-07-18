package holdingshandler

import (
	"fif-calculator/internal/utils"
	"fif-calculator/views/holdings"
	"net/http"
)

func (h *HoldingsHandler) List(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUID(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	allHoldings, err := h.HoldingsRepository.AllHoldings(userId)

	if err != nil {
		http.Error(w, "Failed to get trades", http.StatusInternalServerError)
	}

	viewModels := convertToHoldingViewModels(allHoldings)

	err = holdings.HoldingsList(r.URL.Path, viewModels).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Could not render start page", http.StatusInternalServerError)
	}
}
