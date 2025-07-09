package tradehandler

import (
	"fif-calculator/internal/model"
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

	trade := model.Trade{
		ID:        tradeID,
		Symbol:    r.FormValue("symbol"),
		BuyDate:   r.FormValue("buyDate"),
		Quantity:  quantity,
		Price:     price,
		Currency:  r.FormValue("currency"),
		Action:    r.FormValue("action"),
		HoldingID: holdingID,
	}

	if err := h.TradeRepository.Update(&trade); err != nil {
		http.Error(w, "Failed to update trade", http.StatusInternalServerError)
		return
	}

	redirect := "/holdings/" + holdingIDParam + "/trades/" + tradeIDParam
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}
