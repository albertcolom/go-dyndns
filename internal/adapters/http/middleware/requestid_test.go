package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRequestIdMiddleware(t *testing.T) {
	t.Run("when X-Request-ID header is not present", func(t *testing.T) {
		var receivedRequestID string

		handler := RequestIdMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedRequestID = RequestIDFromContext(r.Context())
			assert.NotEmpty(t, receivedRequestID)

			_, err := uuid.Parse(receivedRequestID)
			assert.NoError(t, err)

			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.NotEmpty(t, receivedRequestID)
		assert.Equal(t, receivedRequestID, w.Header().Get("X-Request-ID"))
	})

	t.Run("when X-Request-ID header is present", func(t *testing.T) {
		expectedRequestID := "existing-request-id"

		handler := RequestIdMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := RequestIDFromContext(r.Context())
			assert.Equal(t, expectedRequestID, requestID)

			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-Request-ID", expectedRequestID)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, expectedRequestID, w.Header().Get("X-Request-ID"))
	})

	t.Run("when X-Request-ID header is present but empty", func(t *testing.T) {
		var receivedRequestID string

		handler := RequestIdMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedRequestID = RequestIDFromContext(r.Context())
			assert.NotEmpty(t, receivedRequestID)

			_, err := uuid.Parse(receivedRequestID)
			assert.NoError(t, err)

			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-Request-ID", "")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.NotEmpty(t, receivedRequestID)
		assert.Equal(t, receivedRequestID, w.Header().Get("X-Request-ID"))
	})
}
