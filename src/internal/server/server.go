package server

import (
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

	s.r.Get("/auth/login", handlers.HandleLoginPage())
	s.r.Post("/auth/login", handlers.HandleLoginForm(s.deps))
	s.r.Get("/auth/logout", handlers.HandleLogout(s.deps))

	s.r.Group(func(r *router.Router) {
		r.Use(AuthMaybeRequiredMW(s.deps))
		r.Get("/", handlers.HandleIndex())
	})

	s.r.Group(func(r *router.Router) {
		r.Use(MustAuthMW(s.deps))
		r.Get("/measure", handlers.HandleMeasure(s.deps))
	})
}

func Start(s *Server) error {
	slog.Info("Starting server", "port", s.port)

	return http.ListenAndServe(s.port, s.r)
}
