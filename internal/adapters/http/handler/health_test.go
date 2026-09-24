package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-dyndns/internal/ports/mocks"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestLivezHandler(t *testing.T) {
	handler := NewHealthHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/livez", nil)
	resp := httptest.NewRecorder()
	handler.Livez(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Empty(t, resp.Body.String())
}

func TestReadyzHandler(t *testing.T) {
	t.Run("No health checker configured", func(t *testing.T) {
		handler := NewHealthHandler(nil)

		req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
		resp := httptest.NewRecorder()
		handler.Readyz(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("Database reachable", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHealthChecker := mocks.NewMockHealthChecker(ctrl)
		mockHealthChecker.EXPECT().Ping(gomock.Any()).Return(nil)
		handler := NewHealthHandler(mockHealthChecker)

		req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
		resp := httptest.NewRecorder()
		handler.Readyz(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("Database unreachable", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHealthChecker := mocks.NewMockHealthChecker(ctrl)
		mockHealthChecker.EXPECT().Ping(gomock.Any()).Return(errors.New("connection refused"))
		handler := NewHealthHandler(mockHealthChecker)

		req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
		resp := httptest.NewRecorder()
		handler.Readyz(resp, req)

		assert.Equal(t, http.StatusServiceUnavailable, resp.Code)
		assert.JSONEq(t, `{"error":"Database unavailable"}`, resp.Body.String())
	})
}
