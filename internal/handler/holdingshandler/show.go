package holdingshandler

import (
	"fif-calculator/views/holdings"
	"net/http"
)

func (h *HoldingsHandler) Show(w http.ResponseWriter, r *http.Request) {
	err := holdings.ViewHolding(r.URL.Path).Render(r.Context(), w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
