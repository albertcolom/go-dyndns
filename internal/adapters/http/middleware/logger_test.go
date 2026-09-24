package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go-dyndns/pkg/logger"
	"go.uber.org/mock/gomock"
)

func TestLoggerMiddleware(t *testing.T) {
	t.Run("logs request details correctly", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockLogger := logger.NewMockLogger(ctrl)

		mockLogger.EXPECT().Info(
			"HTTP",
			"Request",
			logger.Field{Key: "method", Value: "GET"},
			logger.Field{Key: "path", Value: "/test"},
			logger.Field{Key: "status", Value: 200},
			logger.Field{Key: "client_ip", Value: "1.2.3.4"},
			gomock.Any(),
			logger.Field{Key: "request_id", Value: "test-request-id"},
		)

		handler := RequestIdMiddleware()(LoggerMiddleware(mockLogger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})))

		req := httptest.NewRequest("GET", "/test?foo=bar", nil)
		req.RemoteAddr = "1.2.3.4:12345"
		req.Header.Set("X-Request-ID", "test-request-id")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200 but got %d", w.Code)
		}
	})
}
