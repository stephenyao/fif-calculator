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

	r.Get("/", tradeHandler.List)
	r.Get("/trades/new", tradeHandler.NewForm)
	r.Post("/trades", tradeHandler.Create)
	r.Get("/trades/{id}", tradeHandler.Show)

	http.ListenAndServe(":8080", r)
}
