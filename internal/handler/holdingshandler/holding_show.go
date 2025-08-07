package holdingshandler

import (
	"database/sql"
	"errors"
	"fif-calculator/internal/model"
	"fif-calculator/internal/service/costbasisservice"
	"fif-calculator/internal/utils"
	"fif-calculator/internal/viewmodel"
	"fif-calculator/views/holdings"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (h *HoldingsHandler) Show(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")
	id, err := strconv.Atoi(param)

	if err != nil {
		http.Error(w, "id should be a number", http.StatusBadRequest)
	}

	userId, err := utils.GetUID(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	holding, err := h.HoldingsRepository.GetHolding(id, userId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Holding not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	uid, err := utils.GetUID(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	trades, err := h.TradeRepository.GetByHoldingID(id, uid)

	costBasis := costbasisservice.CostBasisBySymbol(trades)[holding.Symbol]

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	vm := viewmodel.HoldingViewModel{
		ID:              id,
		Name:            holding.Name,
		Symbol:          holding.Symbol,
		Currency:        holding.Currency,
		TotalTrades:     fmt.Sprintf("%d", costBasis.TotalTrades),
		CurrentQuantity: strconv.FormatFloat(costBasis.CurrentQuantity, 'f', -1, 64),
		Trades:          convertTradesToViewModel(param, trades),
		CostBasis:       fmt.Sprintf("$%.2f", costBasis.CostBasisNZD),
	}

	err = holdings.ViewHolding(r.URL.Path, vm).Render(r.Context(), w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func convertTradesToViewModel(holdingID string, trades []model.Trade) []viewmodel.TradeViewModel {
	viewTrades := make([]viewmodel.TradeViewModel, len(trades))

	for i, trade := range trades {
		viewTrades[i] = viewmodel.TradeViewModel{
			TransactionDate: trade.BuyDate,
			Quantity:        trade.Quantity,
			Price:           trade.Price,
			Currency:        trade.Currency,
			Action:          trade.Action,
			URL:             "/holdings/" + holdingID + "/trades/" + strconv.Itoa(trade.ID),
			BackURL:         "/holdings/" + holdingID,
		}
	}

	return viewTrades
}
