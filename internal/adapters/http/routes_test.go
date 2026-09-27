package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"go-dyndns/internal/adapters/http/handler"
	"go-dyndns/internal/ports"
	"go-dyndns/internal/ports/mocks"
)

const testToken = "test-token"

func newTestRouter(t *testing.T) (*mocks.MockDNSService, http.Handler) {
	t.Helper()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockService := mocks.NewMockDNSService(ctrl)
	h := handler.NewHandler(mockService)
	healthHandler := handler.NewHealthHandler(nil)

	router := chi.NewRouter()
	RegisterRoutes(router, h, healthHandler, testToken)

	return mockService, router
}

// TestRegisterRoutes_Domains locks down the /v1/domains resource shape: PUT
// updates a domain, GET .../update does the same over GET (routers/DDNS
// clients often can't send PUT), plain GET fetches the record, and DELETE
// removes it.
func TestRegisterRoutes_Domains(t *testing.T) {
	t.Run("PUT updates the domain", func(t *testing.T) {
		mockService, router := newTestRouter(t)
		mockService.EXPECT().Update(gomock.Any(), "example.com", "192.168.1.1").Return(nil)

		req := httptest.NewRequest(http.MethodPut, "/v1/domains/example.com?ip=192.168.1.1&token="+testToken, nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("GET .../update also updates the domain", func(t *testing.T) {
		mockService, router := newTestRouter(t)
		mockService.EXPECT().Update(gomock.Any(), "example.com", "192.168.1.1").Return(nil)

		req := httptest.NewRequest(http.MethodGet, "/v1/domains/example.com/update?ip=192.168.1.1&token="+testToken, nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("GET fetches the domain record", func(t *testing.T) {
		mockService, router := newTestRouter(t)
		mockService.EXPECT().Find(gomock.Any(), "example.com").Return(&ports.Dns{Domain: "example.com"}, nil)

		req := httptest.NewRequest(http.MethodGet, "/v1/domains/example.com?token="+testToken, nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("DELETE removes the domain", func(t *testing.T) {
		mockService, router := newTestRouter(t)
		mockService.EXPECT().Delete(gomock.Any(), "example.com").Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/v1/domains/example.com?token="+testToken, nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("Unauthenticated request is rejected", func(t *testing.T) {
		_, router := newTestRouter(t)

		req := httptest.NewRequest(http.MethodGet, "/v1/domains/example.com", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})
}
