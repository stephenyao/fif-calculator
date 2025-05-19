package handler

import (
	"fif-clacultor/views/trades"
	"net/http"
)

func (h *TradeHandler) NewForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	trades.NewTradeForm().Render(r.Context(), w)
}
