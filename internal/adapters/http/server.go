package http

import (
	"context"
	"errors"
	"net/http"
	"time"

	"go-dyndns/internal/adapters/http/handler"
	"go-dyndns/internal/ports"
)

const (
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 10 * time.Second
	writeTimeout      = 10 * time.Second
	idleTimeout       = 60 * time.Second
)

type Server struct {
	HttpServer *http.Server
}

func NewHTTPServer(h *handler.Handler, healthHandler *handler.HealthHandler, addr, token string, log ports.Logger) *Server {
	router := NewRouter(h, healthHandler, token, log)

	httpServer := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
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
