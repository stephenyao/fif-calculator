package handler

import (
	"fif-clacultor/internal/model"
	"fif-clacultor/internal/repository"
	"fif-clacultor/views/trades"
	"fmt"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
)

type TradeHandler struct {
	Repo *repository.TradeRepository
}

func NewTradeHandler(db *sqlx.DB) *TradeHandler {
	return &TradeHandler{Repo: repository.NewTradeRepository(db)}
}

func (h *TradeHandler) NewForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	trades.NewTradeForm().Render(r.Context(), w)
}

func (h *TradeHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	quantity, _ := strconv.ParseFloat(r.FormValue("quantity"), 64)
	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)

	trade := model.Trade{
		Symbol:   r.FormValue("symbol"),
		BuyDate:  r.FormValue("buyDate"),
		Quantity: quantity,
		Price:    price,
		Currency: r.FormValue("currency"),
		Action:   r.FormValue("action"),
	}

	if err := h.Repo.Insert(&trade); err != nil {
		http.Error(w, "Failed to save trade", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/trades/new", http.StatusSeeOther)
}

func (h *TradeHandler) List(w http.ResponseWriter, r *http.Request) {
	tradeList, err := h.Repo.GetAll()
	fmt.Println(err)
	if err != nil {
		http.Error(w, "Failed to get trades", http.StatusInternalServerError)
	}

	trades.TradeList(tradeList).Render(r.Context(), w)
}
