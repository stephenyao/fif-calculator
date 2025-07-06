package holdingshandler

import (
	"database/sql"
	"errors"
	"fif-calculator/internal/viewmodel"
	"fif-calculator/views/holdings"
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

	holding, err := h.Repo.GetHolding(id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Holding not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	vm := viewmodel.HoldingViewModel{
		ID:       id,
		Name:     holding.Name,
		Symbol:   holding.Symbol,
		Currency: holding.Currency,
	}

	err = holdings.ViewHolding(r.URL.Path, vm).Render(r.Context(), w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
