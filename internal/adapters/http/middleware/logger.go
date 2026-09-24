package middleware

import (
	"net"
	"net/http"
	"strings"
	"time"

	"go-dyndns/internal/ports"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// LoggerMiddleware must be registered after RequestIdMiddleware so the
// request ID it reads from context has already been set.
func LoggerMiddleware(log ports.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rec, r)

			log.Info(r.Context(), "Request",
				"component", "HTTP",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rec.status,
				"client_ip", clientIP(r),
				"duration", time.Since(start),
				"request_id", RequestIDFromContext(r.Context()),
			)
		})
	}
}

func clientIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		return strings.TrimSpace(strings.Split(fwd, ",")[0])
	}

	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}

	return r.RemoteAddr
}
