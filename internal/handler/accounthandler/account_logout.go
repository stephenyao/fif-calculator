package accounthandler

import (
	"net/http"
	"os"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	isLocal := os.Getenv("ENV") == "local"

	cookie := &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false, // Safari allows this on localhost
	}

	if !isLocal {
		cookie.SameSite = http.SameSiteNoneMode
		cookie.Secure = true
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
