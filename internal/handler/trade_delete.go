package handler

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

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
