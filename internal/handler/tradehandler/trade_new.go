package tradehandler

import (
	"database/sql"
	"errors"
	"fif-calculator/internal/utils"
	"fif-calculator/internal/viewmodel"
	"fif-calculator/views/trades"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (h *TradeHandler) New(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		http.Error(w, "Could not parse id", http.StatusBadRequest)
		return
	}

	userId, err := utils.GetUID(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	holding, err := h.HoldingRepository.GetHolding(id, userId)
	actionURL := "/holdings/" + strconv.Itoa(id) + "/trades/new"
	backLink := "/holdings/" + strconv.Itoa(id)
	title := "New trade for " + holding.Name
	vm := viewmodel.CreateTradeFormViewModel(title, actionURL, "", backLink)

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
