package handler

import (
	"context"
	"net/http"
	"sync"

	"golang.org/x/sync/errgroup"

	"go-dyndns/internal/ports"
)

type HealthHandler struct {
	checkers map[string]ports.HealthChecker
}

func NewHealthHandler(checkers map[string]ports.HealthChecker) *HealthHandler {
	return &HealthHandler{checkers: checkers}
}

func (h *HealthHandler) Livez(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *HealthHandler) Readyz(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), defaultTimeout)
	defer cancel()

	var mu sync.Mutex
	checks := make(map[string]string, len(h.checkers))

	g, gctx := errgroup.WithContext(ctx)
	for name, checker := range h.checkers {
		if checker == nil {
			continue
		}

		g.Go(func() error {
			status := "ok"
			if err := checker.Ping(gctx); err != nil {
				status = err.Error()
			}

			mu.Lock()
			checks[name] = status
			mu.Unlock()

			return nil
		})
	}
	_ = g.Wait()

	healthy := true
	for _, status := range checks {
		if status != "ok" {
			healthy = false
			break
		}
	}

	status, code := "ok", http.StatusOK
	if !healthy {
		status, code = "unavailable", http.StatusServiceUnavailable
	}

	JSON(w, code, map[string]any{"status": status, "checks": checks})
}
