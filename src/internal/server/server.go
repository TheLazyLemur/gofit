package server

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/TheLazyLemur/gofit/src/internal/db"
	"github.com/TheLazyLemur/gofit/src/internal/handlers"
	"github.com/TheLazyLemur/gofit/src/router"
)

type dependencies interface {
	DBC() *sql.DB
	Querier() db.Querier
	VersionChecker() string
}

type Server struct {
	port string
	r    *router.Router
	deps dependencies
}

func AuthMaybeRequiredMW(deps dependencies) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := r.Cookie("token")
			if err != nil && err == http.ErrNoCookie {
				ctx := context.WithValue(r.Context(), "user", nil)
				h.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			res, err := deps.Querier().JoinSessionByUserId(r.Context(), deps.DBC(), token.Value)
			if err != nil {
				slog.Error("Error getting user", "err", err)
				ctx := context.WithValue(r.Context(), "user", nil)
				h.ServeHTTP(w, r.WithContext(ctx))
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

func NewServer(port string, deps dependencies) *Server {
	if port[0] != ':' {
		port = ":" + port
	}

	r := router.NewRouter()
	return &Server{
		port: port,
		r:    r,
		deps: deps,
	}
}

func MountRoutes(s *Server) {
	s.r.Get("/health", handlers.HandleHealthCheck(s.deps))
	s.r.Get("/auth/signup", handlers.HandleSignupPage())
	s.r.Post("/auth/signup", handlers.HandleSignupForm(s.deps))

	s.r.Group(func(r *router.Router) {
		r.Use(AuthMaybeRequiredMW(s.deps))
		r.Get("/", handlers.HandleIndex())
	})
}

func Start(s *Server) error {
	slog.Info("Starting server", "port", s.port)

	return http.ListenAndServe(s.port, s.r)
}
