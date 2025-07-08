package holdingshandler

import (
	"fif-calculator/internal/model"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (h *HoldingsHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	name := r.FormValue("name")
	ticker := r.FormValue("ticker")
	currency := r.FormValue("currency")

	record := model.HoldingRecord{
		ID:       id,
		Symbol:   ticker,
		Name:     name,
		Currency: currency,
	}

	if err := h.HoldingsRepository.Update(&record); err != nil {
		http.Error(w, "Failed to update trade", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/holdings/"+strconv.Itoa(id), http.StatusSeeOther)
}
