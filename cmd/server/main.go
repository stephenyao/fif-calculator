package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"

	authmiddleware "fif-calculator/internal/authmiddleware"
	"fif-calculator/internal/handler/accounthandler"
	"fif-calculator/internal/handler/authhandler"
	"fif-calculator/internal/handler/fifhandler"
	"fif-calculator/internal/handler/holdingshandler"
	"fif-calculator/internal/handler/tradehandler"
	"fif-calculator/internal/repository"

	firebase "firebase.google.com/go/v4"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3" // register sqlite3 driver
	"google.golang.org/api/option"
)

func main() {
	// DB
	db := sqlx.MustConnect("sqlite3", "./trades.db")
	repository.InitSchema(db)

	// Firebase
	firebaseApp, err := initFirebaseApp()
	if err != nil {
		log.Fatalf("Failed to init Firebase: %v", err)
	}
	globalAuthClient, err := firebaseApp.Auth(context.Background())
	if err != nil {
		log.Fatalf("Failed to create Firebase Auth client: %v", err)
	}

	// Handlers
	tradeHandler := tradehandler.NewTradeHandler(db)
	fifHandler := fifhandler.NewFIFHandler(db)
	holdingsHandler := holdingshandler.NewHoldingsHandler(db)
	authHandler := authhandler.NewAuthHandler(firebaseApp)
	accountHandler := accounthandler.NewAccountHandler()

	r := chi.NewRouter()

	// ----------------------------
	// STATIC: do NOT apply StripSlashes here
	// ----------------------------
	// Redirect /static -> /static/
	r.Get("/static", http.RedirectHandler("/static/", http.StatusMovedPermanently).ServeHTTP)
	// Serve /static/* from ./static
	staticFS := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", staticFS))

	// ----------------------------
	// PUBLIC ROUTES (can use StripSlashes)
	// ----------------------------
	r.Group(func(pub chi.Router) {
		pub.Use(middleware.StripSlashes)
		pub.Get("/login", authHandler.ShowLoginPage)
		pub.Post("/session-login", authHandler.PostLogin)
	})

	// ----------------------------
	// PROTECTED ROUTES (use StripSlashes + auth)
	// ----------------------------
	r.Group(func(protected chi.Router) {
		protected.Use(middleware.StripSlashes)
		protected.Use(authmiddleware.RequireAuth(globalAuthClient))

		protected.Get("/", holdingsHandler.Index)
		protected.Get("/holdings", holdingsHandler.List)

		protected.Get("/fif", fifHandler.Index)
		protected.Get("/fif/start", fifHandler.New)
		protected.Get("/fif/holding/{id}/year/{year}", fifHandler.GetHolding)
		protected.Post("/fif/start", fifHandler.HoldingsInfo)
		protected.Post("/fif/calculate", fifHandler.ShowCalculation)

		protected.Get("/holdings/new", holdingsHandler.New)
		protected.Post("/holdings/new", holdingsHandler.CreateHolding)
		protected.Get("/holdings/{id}", holdingsHandler.Show)
		protected.Get("/holdings/{id}/trades", holdingsHandler.GetHoldingTrades)
		protected.Post("/holdings/{id}/delete", holdingsHandler.Delete)
		protected.Get("/holdings/{id}/edit", holdingsHandler.EditForm)
		protected.Post("/holdings/{id}/edit", holdingsHandler.Update)

		protected.Get("/holdings/{holdingID}/trades/new", tradeHandler.New)
		protected.Post("/holdings/{holdingID}/trades/new", tradeHandler.Create)
		protected.Get("/holdings/{holdingID}/trades/{tradeID}", tradeHandler.Show)
		protected.Get("/holdings/{holdingID}/trades/{tradeID}/edit", tradeHandler.EditForm)
		protected.Post("/holdings/{holdingID}/trades/{tradeID}/edit", tradeHandler.Update)
		protected.Post("/holdings/{holdingID}/trades/{tradeID}/delete", tradeHandler.Delete)

		protected.Get("/account", accountHandler.Show)
		protected.Post("/logout", accounthandler.Logout)
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}

func initFirebaseApp() (*firebase.App, error) {
	// Load .env for local/dev
	_ = godotenv.Load()

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
