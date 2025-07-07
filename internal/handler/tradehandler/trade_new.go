package tradehandler

import (
	"database/sql"
	"errors"
	"fif-calculator/internal/viewmodel"
	"fif-calculator/views/trades"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (h *TradeHandler) NewForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	vm := viewmodel.CreateTradeFormViewModel("", "/trades")
	err := trades.NewTradeForm(r.URL.Path, vm).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Could not render start page", http.StatusInternalServerError)
	}
}

func (h *TradeHandler) New(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		http.Error(w, "Could not parse id", http.StatusBadRequest)
		return
	}

	holding, err := h.HoldingRepository.GetHolding(id)
	actionURL := "/holdings" + strconv.Itoa(id) + "/trades/new"
	vm := viewmodel.CreateTradeFormViewModel(holding.Name, actionURL)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Could not find holding", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = trades.NewTradeForm(r.URL.Path, vm).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Could not render start page", http.StatusInternalServerError)
		return
	}

}
