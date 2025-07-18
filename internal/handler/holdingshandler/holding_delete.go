package holdingshandler

import (
	"fif-calculator/internal/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (h *HoldingsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		http.Error(w, "Not a valid ID", http.StatusInternalServerError)
	}

	userId, err := utils.GetUID(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.HoldingsRepository.DeleteByID(id, userId)

	if err != nil {
		http.Error(w, "Failed to delete holding", http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/holdings", http.StatusSeeOther)
}
