package main

import (
	"fif-clacultor/views/trades"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/trades/new", http.StatusFound)
	})

	r.Get("/trades/new", func(w http.ResponseWriter, r *http.Request) {
		form := trades.NewTradeForm()
		form.Render(r.Context(), w)
	})

	http.ListenAndServe(":8080", r)
}
