package handler

import (
	"fif-clacultor/views/trades"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (h *TradeHandler) Show(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")
	id, err := strconv.Atoi(param)

	if err != nil {
		http.Error(w, "Not a valid ID", http.StatusInternalServerError)
	}

	trade, err := h.Repo.GetByID(id)
	//println(trade)
	if err != nil {
		http.Error(w, "Failed to get trade", http.StatusInternalServerError)
	}
	trades.TradeDetail(trade).Render(r.Context(), w)
}
