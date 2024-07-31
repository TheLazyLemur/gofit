package handlers

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
				HTMXRedirect(w, r, "/auth/login")
				return
			}

			res, err := deps.Querier().JoinSessionByUserId(r.Context(), deps.DBC(), token.Value)
			if err != nil {
				slog.Error("Error getting user", "err", err)

				ResetTokenCookie(w)
				HTMXRedirect(w, r, "/auth/login")
				return
			}

			serveWithUser(r, w, h, token.Value, res)
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

				ResetTokenCookie(w)
				h.ServeHTTP(w, r)
				return
			}

			serveWithUser(r, w, h, token.Value, res)
		})
	}
}

func serveWithUser(r *http.Request, w http.ResponseWriter, h http.Handler, token string, res db.JoinSessionByUserIdRow) {
	user := db.User{
		ID:        res.UserID,
		Name:      res.Name,
		Email:     res.Email,
		CreatedAt: res.CreatedAt,
	}

	ctx := context.WithValue(r.Context(), "user", user)
	newCtx := context.WithValue(ctx, "token", token)
	h.ServeHTTP(w, r.WithContext(newCtx))
}
