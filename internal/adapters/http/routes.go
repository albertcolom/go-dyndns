package http

import (
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"go-dyndns/internal/adapters/http/handler"
	"go-dyndns/internal/adapters/http/middleware"
	"go-dyndns/internal/ports"
)

func NewRouter(h *handler.Handler, healthHandler *handler.HealthHandler, token string, log ports.Logger) chi.Router {
	router := chi.NewRouter()
	router.Use(chimiddleware.ClientIPFromRemoteAddr)
	router.Use(middleware.RequestIdMiddleware())
	router.Use(middleware.LoggerMiddleware(log))
	router.Use(chimiddleware.Recoverer)

	router.Get("/livez", healthHandler.Livez)
	router.Get("/readyz", healthHandler.Readyz)

	router.Route("/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(token))
			r.Route("/domains", func(r chi.Router) {
				r.Post("/{domain}", h.CreateDomain)
				// GET alias for PUT below: routers/DDNS clients often can only send GET.
				r.Get("/{domain}/update", h.UpdateDomain)
				r.Put("/{domain}", h.UpdateDomain)
				r.Get("/{domain}", h.GetDomain)
				r.Delete("/{domain}", h.DeleteDomain)
			})
		})
	})

	return router
}
