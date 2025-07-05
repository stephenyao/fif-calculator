package holdingshandler

import (
	"fif-calculator/internal/model"
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

	holding := model.HoldingRecord{
		Name:     name,
		Symbol:   ticker,
		Currency: currency,
	}

	err = h.Repo.CreateHolding(&holding)

	if err != nil {
		http.Error(w, "Could not create holding", http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/holdings", http.StatusFound)
}
