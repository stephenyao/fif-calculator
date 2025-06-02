package handler

import (
	"fif-calculator/internal/model"
	"fif-calculator/internal/service/costbasisservice"
	. "fif-calculator/internal/viewmodel"
	"fif-calculator/views/trades"
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

	err = trades.TradeList(
		r.URL.Path,
		tradeList,
		costBasisViewModel(tradesByDescDate),
	).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Could not render start page", http.StatusInternalServerError)
	}
}

func costBasisViewModel(trades []model.Trade) CostBasisViewModel {
	costBasisBySymbol := costbasisservice.CostBasisBySymbol(trades)
	totalCostBasis := 0.0

	for _, costBasisSymbol := range costBasisBySymbol {
		totalCostBasis += costBasisSymbol.CostBasis
	}

	isValidForFIF := totalCostBasis >= 50000

	return CostBasisViewModel{
		costBasisBySymbol,
		totalCostBasis,
		isValidForFIF,
	}
}
