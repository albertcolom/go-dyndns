package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoggerMiddleware(t *testing.T) {
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(io.Discard)

	gin.SetMode(gin.TestMode)

	t.Run("logs request details correctly", func(t *testing.T) {
		w := httptest.NewRecorder()
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("RequestID", "test-request-id")
		})
		router.Use(LoggerMiddleware())

		router.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "OK")
		})

		req, _ := http.NewRequest(http.MethodGet, "/test?param=value", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		router.ServeHTTP(w, req)

		logOutput := logBuf.String()

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, logOutput, "[HTTP] GET /test?param=value")
		assert.Contains(t, logOutput, "| 200 |")      // Status code
		assert.Contains(t, logOutput, "192.168.1.1")  // Client IP
		assert.Contains(t, logOutput, "test-request") // Request ID
	})
}
