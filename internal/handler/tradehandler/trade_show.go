package tradehandler

import (
	"fif-calculator/internal/utils"
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

	uid, err := utils.GetUID(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	trade, err := h.TradeRepository.GetByID(id, uid)
	if err != nil {
		http.Error(w, "Failed to get trade", http.StatusInternalServerError)
		return
	}

	backURL := "/holdings/" + holdingIDParam
	editURL := "/holdings/" + holdingIDParam + "/trades/" + tradeIDParam + "/edit"
	deleteURL := "/holdings/" + holdingIDParam + "/trades/" + tradeIDParam + "/delete"
	err = trades.TradeDetail(r.URL.Path, trade, backURL, editURL, deleteURL).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Could not render start page", http.StatusInternalServerError)
	}
}
