package tradehandler

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (h *TradeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	tradeID, err := strconv.Atoi(chi.URLParam(r, "tradeID"))
	holdingIDParam := chi.URLParam(r, "holdingID")

	if err != nil {
		http.Error(w, "Not a valid ID", http.StatusInternalServerError)
	}

	err = h.TradeRepository.DeleteByID(tradeID)

	if err != nil {
		http.Error(w, "Failed to delete trade", http.StatusInternalServerError)
	}

	redirect := "/holdings/" + holdingIDParam
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}
