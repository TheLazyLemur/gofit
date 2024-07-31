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
}

func HandleHealthCheck(d dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := d.Querier().Ping(r.Context(), d.DBC())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write([]byte("OK"))
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

		token, err := ops.LoginUser(r.Context(), d.DBC(), d.Querier(), email, password)
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
		HTMXRedirect(w, r, "/")
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

		token, err := ops.CreateUser(r.Context(), d.DBC(), d.Querier(), username, password, email)
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
		HTMXRedirect(w, r, "/")
	}
}

func HandleLogout(d dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, ok := r.Context().Value("token").(string)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := d.Querier().DeleteSession(r.Context(), d.DBC(), token); err != nil {
			slog.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ResetTokenCookie(w)
		HTMXRedirect(w, r, "/")

	}
}

func HandleMeasure(deps dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := views.Measure().Render(r.Context(), w); err != nil {
			slog.Error(err.Error())
			return
		}
	}
}

func HandleMeasureWeight(deps dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := views.Weight().Render(r.Context(), w); err != nil {
			slog.Error(err.Error())
			return
		}
	}
}

func HandleMeasureWeightForm(deps dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
