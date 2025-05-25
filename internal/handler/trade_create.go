package handler

import (
	"fif-clacultor/internal/constants"
	"fif-clacultor/internal/model"
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
	symbol := r.FormValue("symbol")
	action := r.FormValue("action")

	existingTrades, err := h.Repo.GetBySymbol(symbol)

	if err != nil {
		http.Error(w, "Invalid symbol", http.StatusInternalServerError)
		return
	}

	err = verifyTradeIsValid(existingTrades, quantity, action)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	trade := model.Trade{
		Symbol:   symbol,
		BuyDate:  r.FormValue("buyDate"),
		Quantity: quantity,
		Price:    price,
		Currency: r.FormValue("currency"),
		Action:   action,
	}

	if err := h.Repo.Insert(&trade); err != nil {
		http.Error(w, "Failed to save trade", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func verifyTradeIsValid(existingTrades []*model.Trade, quantity float64, action string) error {
	switch action {
	case constants.Buy:
		return nil
	case constants.Sell:
		var totalQuantity float64 = 0
		for _, trade := range existingTrades {
			if trade.Action == constants.Buy {
				totalQuantity += trade.Quantity
			} else if trade.Action == constants.Sell {
				totalQuantity -= trade.Quantity
			}
		}

		if totalQuantity >= quantity {
			return nil
		} else {
			return CreateError{message: "The quantity sold is more than the number of holdings in portfolio"}
		}
	default:
		return CreateError{"Unknown trade action"} // unknown and expected case, so return false
	}
}
