package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"go-dyndns/internal/adapters/http/handler"
	"go-dyndns/internal/adapters/http/middleware"
	"go-dyndns/internal/ports"
)

type Server struct {
	HttpServer *http.Server
}

func NewHTTPServer(h *handler.Handler, healthHandler *handler.HealthHandler, addr, token string, log ports.Logger) *Server {
	router := chi.NewRouter()
	router.Use(middleware.RequestIdMiddleware())
	router.Use(middleware.LoggerMiddleware(log))
	router.Use(chimiddleware.Recoverer)

	RegisterRoutes(router, h, healthHandler, token)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	return &Server{
		HttpServer: httpServer,
	}
}

func (s *Server) Start() error {
	if err := s.HttpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.HttpServer.Shutdown(ctx)
}
