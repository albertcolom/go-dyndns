package handler

import (
	"context"
	"net/http"
	"time"

	"go-dyndns/internal/ports"
)

var timeout = 2 * time.Second

type HealthHandler struct {
	checkers map[string]ports.HealthChecker
}

type checkResult struct {
	name string
	err  error
}

func NewHealthHandler(checkers map[string]ports.HealthChecker) *HealthHandler {
	return &HealthHandler{checkers: checkers}
}

func (h *HealthHandler) Livez(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *HealthHandler) Readyz(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), timeout)
	defer cancel()

	active := make(map[string]ports.HealthChecker, len(h.checkers))
	for name, checker := range h.checkers {
		if checker != nil {
			active[name] = checker
		}
	}

	results := make(chan checkResult, len(active))
	for name, checker := range active {
		go func(name string, checker ports.HealthChecker) {
			results <- checkResult{name: name, err: checker.Ping(ctx)}
		}(name, checker)
	}

	checks := make(map[string]string, len(active))
	healthy := true
	ctxDone := false
	for range active {
		if ctxDone {
			break
		}
		select {
		case res := <-results:
			if res.err != nil {
				checks[res.name] = res.err.Error()
				healthy = false
			} else {
				checks[res.name] = "ok"
			}
		case <-ctx.Done():
			healthy = false
			ctxDone = true
		}
	}

	for name := range active {
		if _, reported := checks[name]; !reported {
			checks[name] = "Error: " + ctx.Err().Error()
		}
	}

	status, code := "ok", http.StatusOK
	if !healthy {
		status, code = "unavailable", http.StatusServiceUnavailable
	}

	JSON(w, code, map[string]any{"status": status, "checks": checks})
}
