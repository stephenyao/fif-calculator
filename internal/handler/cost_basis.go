package handler

import (
	"fif-clacultor/internal/model"
	"fif-clacultor/internal/repository"
	. "fif-clacultor/internal/viewmodel"
	"fif-clacultor/views/costbasis"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type CostBasisHandler struct {
	Repo repository.TradeRepository
}

func NewCostBasisHandler(db *sqlx.DB) *CostBasisHandler {
	return &CostBasisHandler{
		Repo: repository.NewTradeRepository(db),
	}
}

func (h *CostBasisHandler) Index(w http.ResponseWriter, r *http.Request) {
	tradeList, err := h.Repo.GetAllByAscendingDate()
	if err != nil {
		http.Error(w, "Failed to get trades", http.StatusInternalServerError)
	}

	viewModel := ccostBasisViewModel(tradeList)

	err = costbasis.Index(r.URL.Path, viewModel).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Could not render start page", http.StatusInternalServerError)
	}
}

func ccostBasisViewModel(trades []model.Trade) CostBasisViewModel {
	costBasisBySymbol := model.CostBasisBySymbol(trades)
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
