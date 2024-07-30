package handlers

import "net/http"

func doRedirect(w http.ResponseWriter, r *http.Request, url string) {
	isHTMX := r.Header.Get("HX-Request") == "true"
	if isHTMX {
		w.Header().Set("HX-Redirect", url)
	} else {
		http.Redirect(w, r, url, http.StatusFound)
	}
}
