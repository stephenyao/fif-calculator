package holdingshandler

import (
	"fif-calculator/internal/repository"
	"github.com/jmoiron/sqlx"
)

type HoldingsHandler struct {
	Repo repository.HoldingsRepository
}

func NewHoldingsHandler(db *sqlx.DB) *HoldingsHandler {
	return &HoldingsHandler{Repo: repository.NewHoldingsRepository(db)}
}
