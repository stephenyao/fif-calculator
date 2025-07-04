package holdingshandler

import (
	"fif-calculator/internal/viewmodel"
	"fif-calculator/views/holdings"
	"net/http"
)

func (h *HoldingsHandler) New(w http.ResponseWriter, r *http.Request) {
	vm := viewmodel.NewHoldingFormViewModel()
	err := holdings.NewHolding(r.URL.Path, vm).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Could not render start page", http.StatusInternalServerError)
	}
}
