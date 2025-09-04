package costbasishandler

import (
	"fif-calculator/internal/repository"
	"github.com/jmoiron/sqlx"
)

type CostBasisHandler struct {
	Repo repository.TradeRepository
}

func NewCostBasisHandler(db *sqlx.DB) *CostBasisHandler {
	return &CostBasisHandler{
		Repo: repository.NewTradeRepository(db),
	}
}
