package handler

import (
	"net/http"

	"go-dyndns/internal/ports"
)

type HealthHandler struct {
	healthChecker ports.HealthChecker
}

func NewHealthHandler(healthChecker ports.HealthChecker) *HealthHandler {
	return &HealthHandler{healthChecker: healthChecker}
}

func (h *HealthHandler) Livez(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *HealthHandler) Readyz(w http.ResponseWriter, r *http.Request) {
	if h.healthChecker == nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	if err := h.healthChecker.Ping(r.Context()); err != nil {
		Error(w, http.StatusServiceUnavailable, "Database unavailable")
		return
	}

	w.WriteHeader(http.StatusOK)
}
