package main

import (
	"fif-clacultor/views"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		h := views.Hello("Stephen")
		h.Render(r.Context(), w)
		//fmt.Fprintln(w, "Hello, world!")
	})

	http.ListenAndServe(":8080", r)
}
