package authmiddleware

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"log"
	"net/http"
	"time"
)

func RequireAuth(client *auth.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session")
			if err != nil || cookie.Value == "" {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			if err != nil {
				log.Println("Failed to get Firebase Auth client:", err)
				http.Error(w, "Internal error", http.StatusInternalServerError)
				return
			}
			start := time.Now()

			// Verify session cookie (not ID token)
			token, err := client.VerifySessionCookie(r.Context(), cookie.Value)
			if err != nil {
				log.Println("Invalid session cookie:", err)
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			log.Println("🕒 verify token took:", time.Since(start))

			// Optionally pass UID down the request context
			ctx := context.WithValue(r.Context(), "uid", token.UID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
