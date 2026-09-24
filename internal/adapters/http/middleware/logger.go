package middleware

import (
	"net"
	"net/http"
	"strings"
	"time"

	"go-dyndns/pkg/logger"
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
func LoggerMiddleware(log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rec, r)

			log.Info(
				"HTTP",
				"Request",
				logger.Field{Key: "method", Value: r.Method},
				logger.Field{Key: "path", Value: r.URL.Path},
				logger.Field{Key: "status", Value: rec.status},
				logger.Field{Key: "client_ip", Value: clientIP(r)},
				logger.Field{Key: "duration", Value: time.Since(start)},
				logger.Field{Key: "request_id", Value: RequestIDFromContext(r.Context())},
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
