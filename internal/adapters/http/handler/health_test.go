package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-dyndns/internal/ports"
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
	t.Run("No health checkers configured", func(t *testing.T) {
		handler := NewHealthHandler(nil)

		req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
		resp := httptest.NewRecorder()
		handler.Readyz(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, `{"status":"ok","checks":{}}`, resp.Body.String())
	})

	t.Run("All checkers reachable", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHealthChecker := mocks.NewMockHealthChecker(ctrl)
		mockHealthChecker.EXPECT().Ping(gomock.Any()).Return(nil)
		mockDNSChecker := mocks.NewMockHealthChecker(ctrl)
		mockDNSChecker.EXPECT().Ping(gomock.Any()).Return(nil)
		handler := NewHealthHandler(map[string]ports.HealthChecker{
			"Database":   mockHealthChecker,
			"DNS server": mockDNSChecker,
		})

		req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
		resp := httptest.NewRecorder()
		handler.Readyz(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, `{"status":"ok","checks":{"Database":"ok","DNS server":"ok"}}`, resp.Body.String())
	})

	t.Run("Database unreachable", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHealthChecker := mocks.NewMockHealthChecker(ctrl)
		mockHealthChecker.EXPECT().Ping(gomock.Any()).Return(errors.New("connection refused"))
		mockDNSChecker := mocks.NewMockHealthChecker(ctrl)
		mockDNSChecker.EXPECT().Ping(gomock.Any()).Return(nil)
		handler := NewHealthHandler(map[string]ports.HealthChecker{
			"Database":   mockHealthChecker,
			"DNS server": mockDNSChecker,
		})

		req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
		resp := httptest.NewRecorder()
		handler.Readyz(resp, req)

		assert.Equal(t, http.StatusServiceUnavailable, resp.Code)
		assert.JSONEq(t, `{"status":"unavailable","checks":{"Database":"connection refused","DNS server":"ok"}}`, resp.Body.String())
	})

	t.Run("DNS server unreachable", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDNSChecker := mocks.NewMockHealthChecker(ctrl)
		mockDNSChecker.EXPECT().Ping(gomock.Any()).Return(errors.New("i/o timeout"))
		handler := NewHealthHandler(map[string]ports.HealthChecker{
			"DNS server": mockDNSChecker,
		})

		req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
		resp := httptest.NewRecorder()
		handler.Readyz(resp, req)

		assert.Equal(t, http.StatusServiceUnavailable, resp.Code)
		assert.JSONEq(t, `{"status":"unavailable","checks":{"DNS server":"i/o timeout"}}`, resp.Body.String())
	})

	t.Run("Slow checker times out", func(t *testing.T) {
		origTimeout := timeout
		timeout = 20 * time.Millisecond
		defer func() { timeout = origTimeout }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSlowChecker := mocks.NewMockHealthChecker(ctrl)
		mockSlowChecker.EXPECT().Ping(gomock.Any()).DoAndReturn(func(ctx context.Context) error {
			<-ctx.Done() // never resolves on its own; only the deadline unblocks it
			return ctx.Err()
		})
		handler := NewHealthHandler(map[string]ports.HealthChecker{
			"Slow": mockSlowChecker,
		})

		req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
		resp := httptest.NewRecorder()
		handler.Readyz(resp, req)

		assert.Equal(t, http.StatusServiceUnavailable, resp.Code)
		assert.JSONEq(t, `{"status":"unavailable","checks":{"Slow":"Error: context deadline exceeded"}}`, resp.Body.String())
	})

	t.Run("Nil checker entry is skipped", func(t *testing.T) {
		handler := NewHealthHandler(map[string]ports.HealthChecker{
			"Database": nil,
		})

		req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
		resp := httptest.NewRecorder()
		handler.Readyz(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, `{"status":"ok","checks":{}}`, resp.Body.String())
	})
}
