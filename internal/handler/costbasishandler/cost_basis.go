package costbasishandler

import (
	"fif-calculator/internal/model"
	"fif-calculator/internal/repository"
	"fif-calculator/internal/service/costbasisservice"
	"fif-calculator/internal/utils"
	. "fif-calculator/internal/viewmodel"
	"fif-calculator/views/costbasis"
	"github.com/jmoiron/sqlx"
	"net/http"
	"time"
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
	userID, err := utils.GetUID(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tradeList, err := h.Repo.GetAllByAscendingDate(userID)
	if err != nil {
		http.Error(w, "Failed to get trades", http.StatusInternalServerError)
		return
	}

	viewModel := ccostBasisViewModel(tradeList)

	err = costbasis.Index(r.URL.Path, viewModel).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Could not render start page", http.StatusInternalServerError)
	}
}

func ccostBasisViewModel(trades []model.Trade) CostBasisViewModel {
	now := time.Now()
	costBasisBySymbol := costbasisservice.CostBasisBySymbol(trades, now)
	totalCostBasis := 0.0

	for _, costBasisSymbol := range costBasisBySymbol {
		totalCostBasis += costBasisSymbol.CostBasisFX
	}

	isValidForFIF := totalCostBasis >= 50000

	return CostBasisViewModel{
		costBasisBySymbol,
		totalCostBasis,
		isValidForFIF,
	}
}
