package tradehandler

import (
	"fif-calculator/views/trades"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (h *TradeHandler) Show(w http.ResponseWriter, r *http.Request) {
	holdingIDParam := chi.URLParam(r, "holdingID")
	tradeIDParam := chi.URLParam(r, "tradeID")
	id, err := strconv.Atoi(tradeIDParam)

	if err != nil {
		http.Error(w, "Not a valid ID", http.StatusInternalServerError)
	}

	trade, err := h.TradeRepository.GetByID(id)
	if err != nil {
		http.Error(w, "Failed to get trade", http.StatusInternalServerError)
	}

	backURL := "/holdings/" + holdingIDParam
	editURL := "/holdings/" + holdingIDParam + "/trades/" + tradeIDParam + "/edit"
	err = trades.TradeDetail(r.URL.Path, trade, backURL, editURL).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Could not render start page", http.StatusInternalServerError)
	}
}
