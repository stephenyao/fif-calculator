package main

import (
	"context"
	"encoding/base64"
	authmiddleware "fif-calculator/internal/authmiddleware"
	"fif-calculator/internal/handler/accounthandler"
	"fif-calculator/internal/handler/authhandler"
	"fif-calculator/internal/handler/costbasishandler"
	"fif-calculator/internal/handler/fifhandler"
	"fif-calculator/internal/handler/holdingshandler"
	"fif-calculator/internal/handler/tradehandler"
	"fif-calculator/internal/repository"
	"fmt"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go/v4"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3" // register sqlite3 driver
	"google.golang.org/api/option"
)

func main() {
	db := sqlx.MustConnect("sqlite3", "./trades.db")
	repository.InitSchema(db)

	r := chi.NewRouter()

	firebaseApp, err := initFirebaseApp()
	globalAuthClient, _ := firebaseApp.Auth(context.Background())

	if err != nil {
		log.Fatalf("Failed to init Firebase: %v", err)
	}

	tradeHandler := tradehandler.NewTradeHandler(db)
	costBasisHandler := costbasishandler.NewCostBasisHandler(db)
	fifHandler := fifhandler.NewFIFHandler(db)
	holdingsHandler := holdingshandler.NewHoldingsHandler(db)
	authHandler := authhandler.NewAuthHandler(firebaseApp)
	accountHandler := accounthandler.NewAccountHandler()

	r.Use(middleware.StripSlashes)

	r.Group(func(protected chi.Router) {
		protected.Use(authmiddleware.RequireAuth(globalAuthClient))
		protected.Get("/holdings", holdingsHandler.List)
		protected.Get("/", holdingsHandler.Index)

		protected.Get("/cost-basis", costBasisHandler.Index)
		protected.Get("/fif", fifHandler.Index)
		protected.Get("/fif/start", fifHandler.New)
		protected.Post("/fif/start", fifHandler.HoldingsInfo)
		protected.Post("/fif/calculate", fifHandler.FIFFormSubmit)
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

		protected.Get("/account", accountHandler.Show)
		protected.Post("/logout", accounthandler.Logout)
	})

	r.Get("/login", authHandler.ShowLoginPage)
	r.Post("/session-login", authHandler.PostLogin)
	log.Fatal(http.ListenAndServe(":8080", r))
}

func initFirebaseApp() (*firebase.App, error) {
	// Load .env (only useful for local dev; no-op in App Platform)
	_ = godotenv.Load()
	fmt.Println("INITILIASING FIREBASE APP")
	b64 := os.Getenv("FIREBASE_KEY_B64")
	if b64 == "" {
		return nil, fmt.Errorf("FIREBASE_KEY_B64 is not set")
	}

	decoded, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode FIREBASE_KEY_B64: %w", err)
	}

	app, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsJSON(decoded))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase App: %w", err)
	}

	return app, nil
}
