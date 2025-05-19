package handler

import (
	"fif-clacultor/internal/model"
	"net/http"
	"strconv"
)

func (h *TradeHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	quantity, _ := strconv.ParseFloat(r.FormValue("quantity"), 64)
	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)

	trade := model.Trade{
		Symbol:   r.FormValue("symbol"),
		BuyDate:  r.FormValue("buyDate"),
		Quantity: quantity,
		Price:    price,
		Currency: r.FormValue("currency"),
		Action:   r.FormValue("action"),
	}

	if err := h.Repo.Insert(&trade); err != nil {
		http.Error(w, "Failed to save trade", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
