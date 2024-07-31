package server

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/TheLazyLemur/gofit"
	"github.com/TheLazyLemur/gofit/src/internal/db"
	"github.com/TheLazyLemur/gofit/src/internal/handlers"
	"github.com/go-chi/chi/v5"
)

type dependencies interface {
	DBC() *sql.DB
	Querier() db.Querier
}

type Server struct {
	port string
	r    *chi.Mux
	deps dependencies
}

func NewServer(port string, deps dependencies) *Server {
	if port[0] != ':' {
		port = ":" + port
	}

	r := chi.NewRouter()
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

	s.r.Group(func(r chi.Router) {
		r.Use(handlers.AuthMaybeRequiredMW(s.deps))
		r.Get("/", handlers.HandleIndex())
	})

	s.r.Group(func(r chi.Router) {
		r.Use(handlers.MustAuthMW(s.deps))
		r.Get("/measure/", handlers.HandleMeasure(s.deps))
		r.Get("/measure/weight", handlers.HandleMeasureWeight(s.deps))
		r.Post("/measure/weight", handlers.HandleMeasureWeightForm(s.deps))
		r.Get("/auth/logout", handlers.HandleLogout(s.deps))
	})

	staticFS := http.FS(gofit.Static)
	s.r.Handle("GET /static/*", http.FileServer(staticFS))
}

func Start(s *Server) error {
	slog.Info("Starting server", "port", s.port)
	return http.ListenAndServe(s.port, s.r)
}
