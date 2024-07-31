package server

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/TheLazyLemur/gofit/src/internal/db"
)

func MustAuthMW(deps dependencies) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := r.Cookie("token")
			if err != nil && err == http.ErrNoCookie {
				isHTMX := r.Header.Get("HX-Request") == "true"
				if isHTMX {
					w.Header().Set("HX-Redirect", "/auth/login")
				} else {
					http.Redirect(w, r, "/auth/login", http.StatusFound)
				}
				return
			}

			res, err := deps.Querier().JoinSessionByUserId(r.Context(), deps.DBC(), token.Value)
			if err != nil {
				slog.Error("Error getting user", "err", err)
				isHTMX := r.Header.Get("HX-Request") == "true"
				if isHTMX {
					w.Header().Set("HX-Redirect", "/auth/login")
				} else {
					http.Redirect(w, r, "/auth/login", http.StatusFound)
				}
				return
			}

			serverWithUser(r, w, h, token.Value, res)
		})
	}
}

func AuthMaybeRequiredMW(deps dependencies) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := r.Cookie("token")
			if err != nil && err == http.ErrNoCookie {
				h.ServeHTTP(w, r)
				return
			}

			res, err := deps.Querier().JoinSessionByUserId(r.Context(), deps.DBC(), token.Value)
			if err != nil {
				slog.Error("Error getting user", "err", err)
				h.ServeHTTP(w, r)
				return
			}

			serverWithUser(r, w, h, token.Value, res)
		})
	}
}

func serverWithUser(r *http.Request, w http.ResponseWriter, h http.Handler, token string, res db.JoinSessionByUserIdRow) {
	user := db.User{
		ID:        res.ID,
		Name:      res.Name,
		Email:     res.Email,
		CreatedAt: res.CreatedAt,
	}

	ctx := context.WithValue(r.Context(), "user", user)
	newCtx := context.WithValue(ctx, "token", token)
	h.ServeHTTP(w, r.WithContext(newCtx))
}
