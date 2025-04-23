package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRequestIdMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("when X-Request-ID header is not present", func(t *testing.T) {
		w := httptest.NewRecorder()
		router := gin.New()
		router.Use(RequestIdMiddleware())

		router.GET("/test", func(c *gin.Context) {
			requestID, exists := c.Get("RequestID")
			assert.True(t, exists)

			_, err := uuid.Parse(requestID.(string))
			assert.NoError(t, err)

			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		router.ServeHTTP(w, req)

		responseHeader := w.Header().Get("X-Request-ID")
		assert.NotEmpty(t, responseHeader)

		_, err := uuid.Parse(responseHeader)
		assert.NoError(t, err)
	})

	t.Run("when X-Request-ID header is present", func(t *testing.T) {
		w := httptest.NewRecorder()
		router := gin.New()
		router.Use(RequestIdMiddleware())

		expectedRequestID := "existing-request-id"

		router.GET("/test", func(c *gin.Context) {
			requestID, exists := c.Get("RequestID")
			assert.True(t, exists)
			assert.Equal(t, expectedRequestID, requestID)

			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-Request-ID", expectedRequestID)
		router.ServeHTTP(w, req)

		assert.Equal(t, expectedRequestID, w.Header().Get("X-Request-ID"))
	})

	t.Run("when X-Request-ID header is present but empty", func(t *testing.T) {
		w := httptest.NewRecorder()
		router := gin.New()
		router.Use(RequestIdMiddleware())

		router.GET("/test", func(c *gin.Context) {
			requestID, exists := c.Get("RequestID")
			assert.True(t, exists)

			_, err := uuid.Parse(requestID.(string))
			assert.NoError(t, err)

			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-Request-ID", "")
		router.ServeHTTP(w, req)

		responseHeader := w.Header().Get("X-Request-ID")
		assert.NotEmpty(t, responseHeader)

		_, err := uuid.Parse(responseHeader)
		assert.NoError(t, err)
	})
}
