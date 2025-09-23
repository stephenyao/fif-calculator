package fifhandler

import (
	"fif-calculator/internal/repository"
	"fif-calculator/internal/service/fifservice"
	"fif-calculator/views/fif"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type FIFHandler struct {
	TradeRepository repository.TradeRepository
	Service         fifservice.FIFService
}

func NewFIFHandler(db *sqlx.DB) *FIFHandler {
	return &FIFHandler{
		TradeRepository: repository.NewTradeRepository(db),
		Service:         fifservice.NewFIFService(repository.NewFIFSQLRepository(db)),
	}
}

func (h *FIFHandler) Index(w http.ResponseWriter, r *http.Request) {
	err := fif.Index(r.URL.Path).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
