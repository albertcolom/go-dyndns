package http

import (
	"github.com/go-chi/chi/v5"

	"go-dyndns/internal/adapters/http/handler"
	"go-dyndns/internal/adapters/http/middleware"
)

// RegisterRoutes wires the API's endpoints onto router.
func RegisterRoutes(router chi.Router, h *handler.Handler, healthHandler *handler.HealthHandler, token string) {
	router.Get("/livez", healthHandler.Livez)
	router.Get("/readyz", healthHandler.Readyz)

	router.Route("/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(token))
			r.Route("/domains", func(r chi.Router) {
				// GET alias for PUT below: routers/DDNS clients often can only send GET.
				r.Get("/{domain}/update", h.UpdateDomain)
				r.Put("/{domain}", h.UpdateDomain)
				r.Get("/{domain}", h.GetDomain)
				r.Delete("/{domain}", h.DeleteDomain)
			})
		})
	})
}
