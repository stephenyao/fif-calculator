package tradehandler

import (
	"fif-calculator/internal/utils"
	"fif-calculator/internal/viewmodel"
	"fif-calculator/views/trades"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (h *TradeHandler) EditForm(w http.ResponseWriter, r *http.Request) {
	holdingIDParam := chi.URLParam(r, "holdingID")
	tradeIDParam := chi.URLParam(r, "tradeID")
	tradeID, err := strconv.Atoi(tradeIDParam)

	if err != nil {
		http.Error(w, "Not a valid ID", http.StatusInternalServerError)
		return
	}

	uid, err := utils.GetUID(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	trade, err := h.TradeRepository.GetByID(tradeID, uid)

	if err != nil {
		http.Error(w, "Failed to get trade", http.StatusInternalServerError)
		return
	}

	title := "Edit trade for " + trade.HoldingName
	vm := viewmodel.CreateTradeFormViewModel(title, "/holdings/"+holdingIDParam+"/trades/"+tradeIDParam+"/edit")

	err = trades.UpdateTradeForm(r.URL.Path, trade, vm).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Could not render start page", http.StatusInternalServerError)
		return
	}
}
