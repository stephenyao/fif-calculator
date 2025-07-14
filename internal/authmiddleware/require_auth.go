package authmiddleware

import (
	"context"
	"log"
	"net/http"

	firebase "firebase.google.com/go/v4"
)

func RequireAuth(firebaseApp *firebase.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session")
			if err != nil || cookie.Value == "" {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			client, err := firebaseApp.Auth(r.Context())
			if err != nil {
				log.Println("Failed to get Firebase Auth client:", err)
				http.Error(w, "Internal error", http.StatusInternalServerError)
				return
			}

			// Verify session cookie (not ID token)
			token, err := client.VerifySessionCookie(r.Context(), cookie.Value)
			if err != nil {
				log.Println("Invalid session cookie:", err)
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			// Optionally pass UID down the request context
			ctx := context.WithValue(r.Context(), "uid", token.UID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
