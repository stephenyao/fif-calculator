package tradehandler

import (
	"fif-calculator/internal/model"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type CreateError struct {
	message string
}

func (c CreateError) Error() string {
	return c.message
}

func (h *TradeHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	quantity, _ := strconv.ParseFloat(r.FormValue("quantity"), 64)
	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	action := r.FormValue("action")
	holdingIDParam := chi.URLParam(r, "id")
	holdingID, err := strconv.Atoi(holdingIDParam)

	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	trade := model.Trade{
		BuyDate:   r.FormValue("buyDate"),
		Quantity:  quantity,
		Price:     price,
		Action:    action,
		HoldingID: holdingID,
	}

	if err := h.TradeRepository.Insert(&trade); err != nil {
		http.Error(w, "Failed to save trade", http.StatusInternalServerError)
		return
	}

	redirect := "/holdings/" + holdingIDParam
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}
