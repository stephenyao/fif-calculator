package handler

import (
	"fif-clacultor/internal/repository"
	"github.com/jmoiron/sqlx"
)

type TradeHandler struct {
	Repo *repository.TradeRepository
}

func NewTradeHandler(db *sqlx.DB) *TradeHandler {
	return &TradeHandler{Repo: repository.NewTradeRepository(db)}
}
