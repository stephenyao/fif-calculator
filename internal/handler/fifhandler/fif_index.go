package fifhandler

import (
	"fif-calculator/internal/repository"
	"fif-calculator/views/fif"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type FIFHandler struct {
	TradeRepository repository.TradeRepository
}

func NewFIFHandler(db *sqlx.DB) *FIFHandler {
	return &FIFHandler{
		TradeRepository: repository.NewTradeRepository(db),
	}
}

func (h *FIFHandler) Index(w http.ResponseWriter, r *http.Request) {
	err := fif.Index(r.URL.Path).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
