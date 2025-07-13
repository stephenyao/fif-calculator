package authhandler

import (
	"context"
	"encoding/json"
	"fif-calculator/views/login"
	firebase "firebase.google.com/go/v4"
	"fmt"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"time"
)

type AuthHandler struct {
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) ShowLoginPage(w http.ResponseWriter, r *http.Request) {
	err := login.Login(r.URL.Path).Render(r.Context(), w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (h *AuthHandler) PostLogin(w http.ResponseWriter, r *http.Request) {
	app, err := firebase.NewApp(
		context.Background(),
		nil,
		option.WithCredentialsFile("private_key.json"),
	)
	if err != nil {
		http.Error(w, "Failed to init Firebase", http.StatusInternalServerError)
		log.Println("Firebase init error:", err)
		return
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		http.Error(w, "Failed to get Firebase auth client", http.StatusInternalServerError)
		log.Println("Firebase auth error:", err)
		return
	}

	var body struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		log.Println("JSON decode error:", err)
		return
	}

	// Verify the ID token before creating a session cookie
	idToken := body.Token
	verifiedToken, err := client.VerifyIDToken(r.Context(), idToken)
	if err != nil {
		http.Error(w, "Invalid ID token", http.StatusUnauthorized)
		log.Println("Token verify error:", err)
		return
	}

	// Set session cookie to expire in 5 days (max 14 days allowed)
	expiresIn := time.Hour * 24 * 5
	sessionCookie, err := client.SessionCookie(r.Context(), idToken, expiresIn)
	if err != nil {
		http.Error(w, "Failed to create session cookie", http.StatusInternalServerError)
		log.Println("Session cookie error:", err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionCookie,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(expiresIn.Seconds()),
	})

	fmt.Fprintf(w, "Logged in as UID: %s", verifiedToken.UID)
}
