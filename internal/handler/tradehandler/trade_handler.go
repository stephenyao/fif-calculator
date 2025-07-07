package tradehandler

import (
	"fif-calculator/internal/repository"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type TradeHandler struct {
	TradeRepository   repository.TradeRepository
	HoldingRepository repository.HoldingsRepository
}

func NewTradeHandler(db *sqlx.DB) *TradeHandler {
	return &TradeHandler{
		TradeRepository:   repository.NewTradeRepository(db),
		HoldingRepository: repository.NewHoldingsRepository(db),
	}
}

func (h *TradeHandler) Index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/trades", http.StatusSeeOther)
}
