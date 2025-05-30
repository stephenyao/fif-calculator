package handler

import (
	"fif-clacultor/views/fif"
	"net/http"
)

func (h *FIFHandler) Start(w http.ResponseWriter, r *http.Request) {
	err := fif.Start(r.URL.Path).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Could not render start fif page", http.StatusInternalServerError)
	}
}
