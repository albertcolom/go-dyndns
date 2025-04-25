package middleware

import (
	"github.com/gin-gonic/gin"
	"go-dyndns/pkg/logger"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoggerMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

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

		router := gin.New()
		router.Use(LoggerMiddleware(mockLogger))

		router.GET("/test", func(c *gin.Context) {
			c.Set("RequestID", "test-request-id")
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest("GET", "/test?foo=bar", nil)
		req.RemoteAddr = "1.2.3.4:12345"
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200 but got %d", w.Code)
		}
	})
}
