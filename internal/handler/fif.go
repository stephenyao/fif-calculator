package handler

import (
	"fif-calculator/internal/repository"
	"fif-calculator/views/fif"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type FIFHandler struct {
	Repo repository.TradeRepository
}

func NewFIFHandler(db *sqlx.DB) *FIFHandler {
	return &FIFHandler{Repo: repository.NewTradeRepository(db)}
}

func (h *FIFHandler) Index(w http.ResponseWriter, r *http.Request) {
	fif.Index(r.URL.Path).Render(r.Context(), w)
}
