package holdingshandler

import (
	"fif-calculator/internal/service/costbasisservice"
	"fif-calculator/internal/utils"
	"fif-calculator/views/holdings"
	"net/http"
	"time"
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
		return
	}

	allTrades, err := h.TradeRepository.GetAllByAscendingDate(userId)

	if err != nil {
		http.Error(w, "Failed to get trades", http.StatusInternalServerError)
		return
	}

	costBasisBySymbol := costbasisservice.CostBasisBySymbol(allTrades, time.Now())

	viewModels := convertToHoldingViewModels(allHoldings, costBasisBySymbol)

	err = holdings.HoldingsList(r.URL.Path, viewModels).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Could not render start page", http.StatusInternalServerError)
	}
}
