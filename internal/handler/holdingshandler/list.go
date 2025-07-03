package holdingshandler

import (
	"fif-calculator/views/holdings"
	"fmt"
	"net/http"
)

func (h *HoldingsHandler) List(w http.ResponseWriter, r *http.Request) {
	allHoldings, err := h.Repo.AllHoldings()

	if err != nil {
		http.Error(w, "Failed to get trades", http.StatusInternalServerError)
	}

	fmt.Printf("AllHoldings: +%v", allHoldings)

	err = holdings.HoldingsList(r.URL.Path).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Could not render start page", http.StatusInternalServerError)
	}
}
