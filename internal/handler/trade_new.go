package handler

import (
	"fif-clacultor/views/trades"
	"net/http"
)

func (h *TradeHandler) NewForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := trades.NewTradeForm(r.URL.Path).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Could not render start page", http.StatusInternalServerError)
	}
}
