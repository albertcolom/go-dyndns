package http

import (
	"context"
	"errors"
	"net/http"

	"go-dyndns/internal/adapters/http/handler"
	"go-dyndns/internal/ports"
)

type Server struct {
	HttpServer *http.Server
}

func NewHTTPServer(h *handler.Handler, healthHandler *handler.HealthHandler, addr, token string, log ports.Logger) *Server {
	router := NewRouter(h, healthHandler, token, log)

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
