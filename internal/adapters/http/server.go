package http

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"go-dyndns/internal/adapters/http/handler"
	"go-dyndns/internal/adapters/http/middleware"
	"go-dyndns/internal/port"
)

type Server struct {
	HttpServer *http.Server
}

func NewHTTPServer(h *handler.Handler, addr, token string, log port.Logger) *Server {
	router := chi.NewRouter()
	// RequestId must run before Logger so the request ID it sets is
	// visible on the request Logger receives.
	router.Use(middleware.RequestIdMiddleware())
	router.Use(middleware.LoggerMiddleware(log))
	router.Use(chimiddleware.Recoverer)

	RegisterRoutes(router, h, token)

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
