package tradehandler

import (
	"fif-calculator/internal/utils"
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

	uid, err := utils.GetUID(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.TradeRepository.DeleteByID(tradeID, uid)

	if err != nil {
		http.Error(w, "Failed to delete trade", http.StatusInternalServerError)
	}

	redirect := "/holdings/" + holdingIDParam
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}
