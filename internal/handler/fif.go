package handler

import (
	"fif-clacultor/views/fif"
	"net/http"
)

type FIFHandler struct {
}

func (h *FIFHandler) Index(w http.ResponseWriter, r *http.Request) {
	fif.Index(r.URL.Path).Render(r.Context(), w)
}
