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

	// dyndns2-compatible endpoint for router/firmware DDNS clients
	router.Route("/nic", func(r chi.Router) {
		r.Use(middleware.BasicAuthMiddleware(token))
		r.Get("/update", h.DynDNS2Update)
	})

	router.Route("/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.TokenAuthMiddleware(token))
			r.Route("/domains", func(r chi.Router) {
				r.Post("/{domain}", h.CreateDomain)
				r.Put("/{domain}", h.UpdateDomain)
				r.Get("/{domain}", h.GetDomain)
				r.Delete("/{domain}", h.DeleteDomain)
			})
		})
	})

	return router
}
