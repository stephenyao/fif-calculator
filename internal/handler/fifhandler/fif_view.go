package fifhandler

import (
	"fif-calculator/views/fif"
	"net/http"
)

func (h *FIFHandler) View(w http.ResponseWriter, r *http.Request) {
	err := fif.View(r.URL.Path).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Could not render start fif page", http.StatusInternalServerError)
	}
}
