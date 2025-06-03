package tradehandler

import (
	"fif-calculator/views/trades"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (h *TradeHandler) EditForm(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")
	id, err := strconv.Atoi(param)

	if err != nil {
		http.Error(w, "Not a valid ID", http.StatusInternalServerError)
	}

	trade, err := h.Repo.GetByID(id)

	if err != nil {
		http.Error(w, "Failed to get trade", http.StatusInternalServerError)
	}

	err = trades.UpdateTradeForm(r.URL.Path, trade).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Could not render start page", http.StatusInternalServerError)
	}
}
