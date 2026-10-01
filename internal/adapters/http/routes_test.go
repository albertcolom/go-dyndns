package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"go-dyndns/internal/adapters/http/handler"
	"go-dyndns/internal/adapters/logger"
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
	log := logger.NewSlogLogger("error") // keep test output quiet

	router := NewRouter(h, healthHandler, testToken, log)

	return mockService, router
}

// TestNewRouter_Domains locks down the /v1/domains resource shape: POST
// creates a domain, PUT updates it, plain GET fetches the record, and
// DELETE removes it.
func TestNewRouter_Domains(t *testing.T) {
	t.Run("POST creates the domain", func(t *testing.T) {
		mockService, router := newTestRouter(t)
		mockService.EXPECT().Create(gomock.Any(), "example.com", "192.168.1.1").Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/v1/domains/example.com?ip=192.168.1.1&token="+testToken, nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
	})

	t.Run("POST rejects a domain that already exists", func(t *testing.T) {
		mockService, router := newTestRouter(t)
		mockService.EXPECT().Create(gomock.Any(), "example.com", "192.168.1.1").Return(ports.ErrDomainExists)

		req := httptest.NewRequest(http.MethodPost, "/v1/domains/example.com?ip=192.168.1.1&token="+testToken, nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusConflict, resp.Code)
	})

	t.Run("PUT updates the domain", func(t *testing.T) {
		mockService, router := newTestRouter(t)
		mockService.EXPECT().Update(gomock.Any(), "example.com", "192.168.1.1").Return(nil)

		req := httptest.NewRequest(http.MethodPut, "/v1/domains/example.com?ip=192.168.1.1&token="+testToken, nil)
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

// TestNewRouter_DynDNS2 locks down the /nic/update dyndns2-compat route used
// by router/firmware DDNS clients (e.g. AVM FritzBox), which authenticate
// via HTTP Basic Auth rather than a bearer token or query param.
func TestNewRouter_DynDNS2(t *testing.T) {
	t.Run("Authenticated request creates the hostname", func(t *testing.T) {
		mockService, router := newTestRouter(t)
		mockService.EXPECT().Find(gomock.Any(), "home.example.com").Return(nil, nil)
		mockService.EXPECT().Create(gomock.Any(), "home.example.com", "192.168.1.1").Return(nil)

		req := httptest.NewRequest(http.MethodGet, "/nic/update?hostname=home.example.com&myip=192.168.1.1", nil)
		req.SetBasicAuth("any-username", testToken)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, "good 192.168.1.1\n", resp.Body.String())
	})

	t.Run("Query token is not accepted, only basic auth", func(t *testing.T) {
		_, router := newTestRouter(t)

		req := httptest.NewRequest(http.MethodGet, "/nic/update?hostname=home.example.com&myip=192.168.1.1&token="+testToken, nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		assert.Equal(t, "badauth\n", resp.Body.String())
	})

	t.Run("Wrong basic auth password is rejected", func(t *testing.T) {
		_, router := newTestRouter(t)

		req := httptest.NewRequest(http.MethodGet, "/nic/update?hostname=home.example.com&myip=192.168.1.1", nil)
		req.SetBasicAuth("any-username", "wrong")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		assert.Equal(t, "badauth\n", resp.Body.String())
	})
}

// TestNewRouter_Middleware confirms the top-level middleware stack (request
// ID, logging, panic recovery) is wired in ahead of the routes it dispatches
// to.
func TestNewRouter_Middleware(t *testing.T) {
	t.Run("Health check bypasses auth", func(t *testing.T) {
		_, router := newTestRouter(t)

		req := httptest.NewRequest(http.MethodGet, "/livez", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("Request ID middleware tags the response", func(t *testing.T) {
		_, router := newTestRouter(t)

		req := httptest.NewRequest(http.MethodGet, "/livez", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.NotEmpty(t, resp.Header().Get("X-Request-ID"))
	})

	t.Run("Recoverer turns a handler panic into a 500", func(t *testing.T) {
		mockService, router := newTestRouter(t)
		mockService.EXPECT().Find(gomock.Any(), "example.com").DoAndReturn(
			func(context.Context, string) (*ports.Dns, error) {
				panic("boom")
			},
		)

		req := httptest.NewRequest(http.MethodGet, "/v1/domains/example.com?token="+testToken, nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("Auto-detects the client IP from the connection when ip is omitted", func(t *testing.T) {
		mockService, router := newTestRouter(t)
		mockService.EXPECT().Update(gomock.Any(), "example.com", "203.0.113.42").Return(nil)

		req := httptest.NewRequest(http.MethodPut, "/v1/domains/example.com?token="+testToken, nil)
		req.RemoteAddr = "203.0.113.42:54321"
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})
}
