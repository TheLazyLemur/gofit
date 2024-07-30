package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/TheLazyLemur/gofit/src/internal/db"
	"github.com/TheLazyLemur/gofit/src/internal/ops"
	"github.com/TheLazyLemur/gofit/src/internal/views"
)

type dependencies interface {
	DBC() *sql.DB
	Querier() db.Querier
	VersionChecker() string
}

func HandleHealthCheck(d dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := d.Querier().Ping(r.Context(), d.DBC())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write([]byte("OK"))
		w.Write([]byte(d.VersionChecker()))
	}
}

func HandleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := views.Home().Render(r.Context(), w); err != nil {
			slog.Error(err.Error())
			return
		}
	}
}

func HandleLoginPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := views.Login().Render(r.Context(), w); err != nil {
			slog.Error(err.Error())
			return
		}
	}
}

func HandleLoginForm(d dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		token, err := ops.LoginUser(r.Context(), d, email, password)
		if err != nil {
			slog.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		cookie := http.Cookie{
			Name:     "token",
			Value:    token,
			HttpOnly: true,
			Path:     "/",
			Expires:  time.Now().Add(time.Hour * 24 * 7),
		}

		http.SetCookie(w, &cookie)
		w.Header().Set("HX-Redirect", "/")
	}
}

func HandleSignupPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := views.Signup().Render(r.Context(), w); err != nil {
			slog.Error(err.Error())
			return
		}
	}
}

func HandleSignupForm(d dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")
		email := r.FormValue("email")

		if username == "" || password == "" || email == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		token, err := ops.CreateUser(r.Context(), d, username, password, email)
		if err != nil {
			slog.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		cookie := http.Cookie{
			Name:     "token",
			Value:    token,
			HttpOnly: true,
			Path:     "/",
			Expires:  time.Now().Add(time.Hour * 24 * 7),
		}

		http.SetCookie(w, &cookie)
		w.Header().Set("HX-Redirect", "/")
	}
}
