package fifhandler

import (
	. "fif-calculator/internal/viewmodel"
	"fif-calculator/views/fif"
	"net/http"
)

func (h *FIFHandler) New(w http.ResponseWriter, r *http.Request) {
	vm := CreateFIFStartViewModel()
	err := fif.Start(r.URL.Path, vm).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Could not render start fif page", http.StatusInternalServerError)
	}
}
