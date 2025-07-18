package accounthandler

import (
	"fif-calculator/views/account"
	"net/http"
)

func (h *AccountHandler) Show(w http.ResponseWriter, r *http.Request) {
	err := account.ShowAccount(r.URL.Path).Render(r.Context(), w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
