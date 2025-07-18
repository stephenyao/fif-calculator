package holdingshandler

import (
	"errors"
	"fif-calculator/internal/model"
	"fif-calculator/internal/utils"
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

func (h *HoldingsHandler) CreateHolding(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		http.Error(w, "Could not parse form", http.StatusInternalServerError)
	}

	name := r.FormValue("name")
	ticker := r.FormValue("ticker")
	currency := r.FormValue("currency")
	userId, err := utils.GetUID(r.Context())

	if errors.Is(err, utils.UIDError{}) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	holding := model.HoldingRecord{
		UserID:   userId,
		Name:     name,
		Symbol:   ticker,
		Currency: currency,
	}

	err = h.HoldingsRepository.CreateHolding(&holding)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/holdings", http.StatusFound)
}
