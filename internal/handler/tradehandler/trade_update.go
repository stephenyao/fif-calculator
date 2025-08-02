package tradehandler

import (
	"fif-calculator/internal/model"
	"fif-calculator/internal/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (h *TradeHandler) Update(w http.ResponseWriter, r *http.Request) {
	tradeIDParam := chi.URLParam(r, "tradeID")
	tradeID, err := strconv.Atoi(tradeIDParam)
	holdingIDParam := chi.URLParam(r, "holdingID")
	holdingID, err := strconv.Atoi(holdingIDParam)

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	quantity, _ := strconv.ParseFloat(r.FormValue("quantity"), 64)
	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	rate, _ := strconv.ParseFloat(r.FormValue("rate"), 64)

	trade := model.Trade{
		ID:           tradeID,
		BuyDate:      r.FormValue("buyDate"),
		Quantity:     quantity,
		Price:        price,
		ExchangeRate: rate,
		Action:       r.FormValue("action"),
		HoldingID:    holdingID,
	}

	uid, err := utils.GetUID(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.TradeRepository.Update(uid, &trade); err != nil {
		http.Error(w, "Failed to update trade", http.StatusInternalServerError)
		return
	}

	redirect := "/holdings/" + holdingIDParam + "/trades/" + tradeIDParam
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}
