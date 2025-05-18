package handler

import (
	"fif-clacultor/internal/model"
	"fif-clacultor/internal/repository"
	"fif-clacultor/views/trades"
	"github.com/go-chi/chi/v5"
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

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *TradeHandler) List(w http.ResponseWriter, r *http.Request) {
	tradeList, err := h.Repo.GetAll()
	if err != nil {
		http.Error(w, "Failed to get trades", http.StatusInternalServerError)
	}

	trades.TradeList(tradeList).Render(r.Context(), w)
}

func (h *TradeHandler) Show(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")
	id, err := strconv.Atoi(param)

	if err != nil {
		http.Error(w, "Not a valid ID", http.StatusInternalServerError)
	}

	trade, err := h.Repo.GetByID(id)
	//println(trade)
	if err != nil {
		http.Error(w, "Failed to get trade", http.StatusInternalServerError)
	}
	trades.TradeDetail(trade).Render(r.Context(), w)
}

func (h *TradeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		http.Error(w, "Not a valid ID", http.StatusInternalServerError)
	}

	err = h.Repo.DeleteByID(id)

	if err != nil {
		http.Error(w, "Failed to delete trade", http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/trades", http.StatusSeeOther)
}

func (h *TradeHandler) EditForm(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")
	id, err := strconv.Atoi(param)

	if err != nil {
		http.Error(w, "Not a valid ID", http.StatusInternalServerError)
	}

	trade, err := h.Repo.GetByID(id)

	if err != nil {
		http.Error(w, "Failed to get trade", http.StatusInternalServerError)
	}

	trades.UpdateTradeForm(trade).Render(r.Context(), w)
}

func (h *TradeHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	quantity, _ := strconv.ParseFloat(r.FormValue("quantity"), 64)
	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)

	trade := model.Trade{
		ID:       id,
		Symbol:   r.FormValue("symbol"),
		BuyDate:  r.FormValue("buyDate"),
		Quantity: quantity,
		Price:    price,
		Currency: r.FormValue("currency"),
		Action:   r.FormValue("action"),
	}

	if err := h.Repo.Update(&trade); err != nil {
		http.Error(w, "Failed to update trade", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/trades/"+strconv.Itoa(id), http.StatusSeeOther)
}
