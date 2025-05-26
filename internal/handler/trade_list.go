package handler

import (
	"fif-clacultor/internal/model"
	. "fif-clacultor/internal/view_model"
	"fif-clacultor/views/trades"
	"net/http"
	"slices"
)

func (h *TradeHandler) List(w http.ResponseWriter, r *http.Request) {
	tradeList, err := h.Repo.GetAll()

	if err != nil {
		http.Error(w, "Failed to get trades", http.StatusInternalServerError)
	}

	tradesByDescDate := slices.Clone(tradeList)
	slices.Reverse(tradesByDescDate)

	trades.TradeList(tradeList, costBasisViewModel(tradesByDescDate)).Render(r.Context(), w)
}

func costBasisViewModel(trades []model.Trade) CostBasisViewModel {
	costBasisBySymbol := model.CostBasisBySymbol(trades)
	totalCostBasis := 0.0

	for _, costBasis := range costBasisBySymbol {
		totalCostBasis += costBasis
	}

	isValidForFIF := totalCostBasis >= 50000

	return CostBasisViewModel{
		costBasisBySymbol,
		totalCostBasis,
		isValidForFIF,
	}
}
