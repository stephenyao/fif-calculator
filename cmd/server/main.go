package main

import (
	"fif-clacultor/internal/handler"
	"fif-clacultor/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // register sqlite3 driver
	"net/http"
)

func main() {
	db := sqlx.MustConnect("sqlite3", "./trades.db")
	repository.InitSchema(db)

	r := chi.NewRouter()

	tradeHandler := handler.NewTradeHandler(db)

	r.Use(middleware.StripSlashes)

	r.Get("/", tradeHandler.List)
	r.Get("/trades/new", tradeHandler.NewForm)
	r.Get("/trades", tradeHandler.List)
	r.Post("/trades", tradeHandler.Create)
	r.Get("/trades/{id}", tradeHandler.Show)
	r.Post("/trades/{id}/delete", tradeHandler.Delete)
	r.Get("/trades/{id}/edit", tradeHandler.EditForm)
	r.Post("/trades/{id}/edit", tradeHandler.Update)

	http.ListenAndServe(":8080", r)
}
