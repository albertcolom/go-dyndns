package http

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"go-dyndns/internal/adapters/http/middleware"
	"go-dyndns/pkg/logger"
)

type Server struct {
	HttpServer *http.Server
}

func NewHTTPServer(handler *Handler, addr, token string, log logger.Logger) *Server {
	router := chi.NewRouter()
	// RequestId must run before Logger so the request ID it sets is
	// visible on the request Logger receives.
	router.Use(middleware.RequestIdMiddleware())
	router.Use(middleware.LoggerMiddleware(log))
	router.Use(chimiddleware.Recoverer)

	router.Route("/v1", func(r chi.Router) {
		r.Get("/health", handler.Health)

		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(token))
			r.Get("/update", handler.UpdateIp)
			r.Get("/get", handler.GetIp)
		})
	})

	httpServer := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	return &Server{
		HttpServer: httpServer,
	}
}

func (s *Server) Start() error {
	return s.HttpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.HttpServer.Shutdown(ctx)
}
