package handlers

import (
	"net/http"
	"time"
)

func HTMXRedirect(w http.ResponseWriter, r *http.Request, url string) {
	isHTMX := r.Header.Get("HX-Request") == "true"
	if isHTMX {
		w.Header().Set("HX-Redirect", url)
	} else {
		http.Redirect(w, r, url, http.StatusFound)
	}
}

func ResetTokenCookie(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:     "token",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * -24 * 7),
	}

	http.SetCookie(w, &cookie)
}
