package main

import (
	"context"
	authmiddleware "fif-calculator/internal/authmiddleware"
	"fif-calculator/internal/handler/authhandler"
	"fif-calculator/internal/handler/costbasishandler"
	"fif-calculator/internal/handler/fifhandler"
	"fif-calculator/internal/handler/holdingshandler"
	"fif-calculator/internal/handler/tradehandler"
	"fif-calculator/internal/repository"
	firebase "firebase.google.com/go/v4"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // register sqlite3 driver
	"google.golang.org/api/option"
	"log"
	"net/http"
)

func main() {
	db := sqlx.MustConnect("sqlite3", "./trades.db")
	repository.InitSchema(db)

	r := chi.NewRouter()

	opt := option.WithCredentialsFile("private_key.json")
	firebaseApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	tradeHandler := tradehandler.NewTradeHandler(db)
	costBasisHandler := costbasishandler.NewCostBasisHandler(db)
	fifHandler := fifhandler.NewFIFHandler(db)
	holdingsHandler := holdingshandler.NewHoldingsHandler(db)
	authHandler := authhandler.NewAuthHandler(firebaseApp)

	r.Use(middleware.StripSlashes)

	r.Group(func(protected chi.Router) {
		protected.Use(authmiddleware.RequireAuth(firebaseApp))
		protected.Get("/holdings", holdingsHandler.List)
		protected.Get("/", holdingsHandler.Index)

		protected.Post("/trades", tradeHandler.Create)
		protected.Get("/trades/{id}", tradeHandler.Show)
		protected.Get("/cost-basis", costBasisHandler.Index)
		protected.Get("/fif", fifHandler.Index)
		protected.Get("/fif/start", fifHandler.New)
		protected.Post("/fif/start", fifHandler.HoldingsInfo)
		protected.Post("/fif/calculate", fifHandler.FIFFormSubmit)
		protected.Get("/fif/view/{id}", fifHandler.View)
		protected.Get("/holdings/new", holdingsHandler.New)
		protected.Post("/holdings/new", holdingsHandler.CreateHolding)
		protected.Get("/holdings/{id}", holdingsHandler.Show)
		protected.Post("/holdings/{id}/delete", holdingsHandler.Delete)
		protected.Get("/holdings/{id}/edit", holdingsHandler.EditForm)
		protected.Post("/holdings/{id}/edit", holdingsHandler.Update)
		protected.Get("/holdings/{id}/trades/new", tradeHandler.New)
		protected.Post("/holdings/{id}/trades/new", tradeHandler.Create)
		protected.Get("/holdings/{holdingID}/trades/{tradeID}", tradeHandler.Show)
		protected.Get("/holdings/{holdingID}/trades/{tradeID}/edit", tradeHandler.EditForm)
		protected.Post("/holdings/{holdingID}/trades/{tradeID}/edit", tradeHandler.Update)
		protected.Post("/holdings/{holdingID}/trades/{tradeID}/delete", tradeHandler.Delete)
	})

	r.Get("/login", authHandler.ShowLoginPage)
	r.Post("/session-login", authHandler.PostLogin)
	log.Fatal(http.ListenAndServe(":8080", r))
}
