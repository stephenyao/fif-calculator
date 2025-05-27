package handler

import (
	"fif-clacultor/views/costbasis"
	"net/http"
)

type CostBasisHandler struct {
}

func (h *CostBasisHandler) Index(w http.ResponseWriter, r *http.Request) {
	costbasis.Index(r.URL.Path).Render(r.Context(), w)
}
