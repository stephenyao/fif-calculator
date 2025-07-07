package main

import (
	"fif-calculator/internal/handler/costbasishandler"
	"fif-calculator/internal/handler/fifhandler"
	"fif-calculator/internal/handler/holdingshandler"
	"fif-calculator/internal/handler/tradehandler"
	"fif-calculator/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // register sqlite3 driver
	"log"
	"net/http"
)

func main() {
	db := sqlx.MustConnect("sqlite3", "./trades.db")
	repository.InitSchema(db)

	r := chi.NewRouter()

	tradeHandler := tradehandler.NewTradeHandler(db)
	costBasisHandler := costbasishandler.NewCostBasisHandler(db)
	fifHandler := fifhandler.NewFIFHandler(db)
	holdingsHandler := holdingshandler.NewHoldingsHandler(db)

	r.Use(middleware.StripSlashes)

	r.Get("/", tradeHandler.Index)
	r.Get("/trades/new", tradeHandler.NewForm)
	r.Get("/trades", tradeHandler.List)
	r.Post("/trades", tradeHandler.Create)
	r.Get("/trades/{id}", tradeHandler.Show)
	r.Post("/trades/{id}/delete", tradeHandler.Delete)
	r.Get("/trades/{id}/edit", tradeHandler.EditForm)
	r.Post("/trades/{id}/edit", tradeHandler.Update)
	r.Get("/cost-basis", costBasisHandler.Index)
	r.Get("/fif", fifHandler.Index)
	r.Get("/fif/start", fifHandler.New)
	r.Post("/fif/start", fifHandler.HoldingsInfo)
	r.Post("/fif/calculate", fifHandler.FIFFormSubmit)
	r.Get("/fif/view/{id}", fifHandler.View)
	r.Get("/holdings", holdingsHandler.List)
	r.Get("/holdings/new", holdingsHandler.New)
	r.Post("/holdings/new", holdingsHandler.CreateHolding)
	r.Get("/holdings/{id}", holdingsHandler.Show)
	r.Post("/holdings/{id}/delete", holdingsHandler.Delete)
	r.Get("/holdings/{id}/edit", holdingsHandler.EditForm)
	r.Post("/holdings/{id}/edit", holdingsHandler.Update)
	r.Get("/holdings/{id}/trades/new", tradeHandler.New)
	r.Post("/holdings/{id}/trades/new", tradeHandler.Create)
	log.Fatal(http.ListenAndServe(":8080", r))
}
