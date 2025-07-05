package holdingshandler

import (
	"fif-calculator/internal/viewmodel"
	"fif-calculator/views/holdings"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (h *HoldingsHandler) EditForm(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")
	id, err := strconv.Atoi(param)

	if err != nil {
		http.Error(w, "Not a valid ID", http.StatusInternalServerError)
	}

	holding, err := h.Repo.GetHolding(id)
	vm := viewmodel.NewHoldingFormViewModelFromRecord(holding)

	if err != nil {
		http.Error(w, "Failed to get trade", http.StatusInternalServerError)
	}

	err = holdings.EditHolding(r.URL.Path, vm).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Could not render start page", http.StatusInternalServerError)
	}
}
