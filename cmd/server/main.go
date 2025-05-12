package main

import (
	"fif-clacultor/internal/handler"
	"fif-clacultor/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // register sqlite3 driver
	"net/http"
)

func main() {
	db := sqlx.MustConnect("sqlite3", "./trades.db")
	repository.InitSchema(db)

	r := chi.NewRouter()

	tradeHandler := handler.NewTradeHandler(db)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/trades/new", http.StatusFound)
	})

	r.Get("/trades/new", tradeHandler.NewForm)
	r.Post("/trades", tradeHandler.Create)

	http.ListenAndServe(":8080", r)
}
