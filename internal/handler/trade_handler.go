package handler

import (
	"fif-clacultor/internal/repository"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type TradeHandler struct {
	Repo repository.TradeRepository
}

func NewTradeHandler(db *sqlx.DB) *TradeHandler {
	return &TradeHandler{Repo: repository.NewTradeRepository(db)}
}

func (h *TradeHandler) Index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/trades", http.StatusSeeOther)
}
