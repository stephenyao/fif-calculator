package handler

import (
	"fif-clacultor/internal/model"
	"fif-clacultor/views/trades"
	"fmt"
	"net/http"
)

func (h *TradeHandler) List(w http.ResponseWriter, r *http.Request) {
	tradeList, err := h.Repo.GetAll()

	if err != nil {
		http.Error(w, "Failed to get trades", http.StatusInternalServerError)
	}

	trades.TradeList(tradeList).Render(r.Context(), w)
}

func computeCostBasisBySymbol(trades []*model.Trade) map[string]float64 {
	tradesBySymbol := make(map[string][]*model.Trade)

	for _, trade := range trades {
		tradesBySymbol[trade.Symbol] = append(tradesBySymbol[trade.Symbol], trade)
	}

	costBasisBySymbol := make(map[string]float64)

	for symbol, trades := range tradesBySymbol {
		var costBasis float64 = 0
		for _, trade := range trades {
			switch trade.Action {
			case "buy":
				costBasis += trade.Price * trade.Quantity
			case "sell":
				costBasis -= trade.Price * trade.Quantity
			}
		}

		costBasisBySymbol[symbol] = costBasis
	}

	for symbol, costBasis := range costBasisBySymbol {
		fmt.Printf("The cost basis of %s is %f", symbol, costBasis)
	}

	return costBasisBySymbol
}
