package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/TheLazyLemur/gofit/src/internal/db"
)

func AuthMaybeRequiredMW(deps dependencies) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("auth maybe required")
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

			user := db.User{
				ID:        res.ID,
				Name:      res.Name,
				Email:     res.Email,
				CreatedAt: res.CreatedAt,
			}
			ctx := context.WithValue(r.Context(), "user", user)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
