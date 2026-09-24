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
			r.Get("/update", h.UpdateIp)
			r.Get("/get", h.GetIp)
		})
	})
}
