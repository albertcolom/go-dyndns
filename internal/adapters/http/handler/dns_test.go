package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"go-dyndns/internal/ports"
	"go-dyndns/internal/ports/mocks"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// newRequestWithDomain builds a request carrying domain as a chi URL param,
// mirroring what the router injects for /v1/domains/{domain}.
func newRequestWithDomain(method, target, domain string) *http.Request {
	req := httptest.NewRequest(method, target, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("domain", domain)
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

func TestCreateDomainHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockDNSService(ctrl)
	handler := NewHandler(mockService)

	domain := "example.com"
	ip := "192.168.1.1"

	t.Run("Create successful", func(t *testing.T) {
		mockService.EXPECT().Create(gomock.Any(), domain, ip).Return(nil)

		req := newRequestWithDomain(http.MethodPost, fmt.Sprintf("/domains/%s?ip=%s", domain, ip), domain)
		resp := httptest.NewRecorder()
		handler.CreateDomain(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
		assert.JSONEq(t, fmt.Sprintf("{\"message\":\"Created %s with %s\"}", domain, ip), resp.Body.String())
	})

	t.Run("Auto-detects IP from the connection when ip parameter is omitted", func(t *testing.T) {
		detectedIP := "203.0.113.42"
		mockService.EXPECT().Create(gomock.Any(), domain, detectedIP).Return(nil)

		req := newRequestWithDomain(http.MethodPost, fmt.Sprintf("/domains/%s", domain), domain)
		req.RemoteAddr = detectedIP + ":54321"
		resp := httptest.NewRecorder()
		handler.CreateDomain(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
		assert.JSONEq(t, fmt.Sprintf("{\"message\":\"Created %s with %s\"}", domain, detectedIP), resp.Body.String())
	})

	t.Run("Failed to determine IP when RemoteAddr has no port", func(t *testing.T) {
		req := newRequestWithDomain(http.MethodPost, fmt.Sprintf("/domains/%s", domain), domain)
		req.RemoteAddr = "not-a-valid-remote-addr"
		resp := httptest.NewRecorder()
		handler.CreateDomain(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, `{"error":"Unable to determine IP"}`, resp.Body.String())
	})

	t.Run("Domain already exists maps to conflict", func(t *testing.T) {
		mockService.EXPECT().Create(gomock.Any(), domain, ip).Return(ports.ErrDomainExists)

		req := newRequestWithDomain(http.MethodPost, fmt.Sprintf("/domains/%s?ip=%s", domain, ip), domain)
		resp := httptest.NewRecorder()
		handler.CreateDomain(resp, req)

		assert.Equal(t, http.StatusConflict, resp.Code)
		assert.JSONEq(t, `{"error":"Domain already exists"}`, resp.Body.String())
	})

	t.Run("Failed validation error maps to bad request", func(t *testing.T) {
		mockService.EXPECT().Create(gomock.Any(), domain, ip).Return(ports.ErrInvalidDomain)

		req := newRequestWithDomain(http.MethodPost, fmt.Sprintf("/domains/%s?ip=%s", domain, ip), domain)
		resp := httptest.NewRecorder()
		handler.CreateDomain(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, fmt.Sprintf(`{"error":%q}`, ports.ErrInvalidDomain.Error()), resp.Body.String())
	})

	t.Run("Failed unexpected error", func(t *testing.T) {
		mockService.EXPECT().Create(gomock.Any(), domain, ip).Return(fmt.Errorf("some error"))

		req := newRequestWithDomain(http.MethodPost, fmt.Sprintf("/domains/%s?ip=%s", domain, ip), domain)
		resp := httptest.NewRecorder()
		handler.CreateDomain(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.JSONEq(t, `{"error":"some error"}`, resp.Body.String())
	})
}

func TestUpdateDomainHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockDNSService(ctrl)
	handler := NewHandler(mockService)

	domain := "example.com"
	ip := "192.168.1.1"

	t.Run("Update successful", func(t *testing.T) {
		mockService.EXPECT().Update(gomock.Any(), domain, ip).Return(nil)

		req := newRequestWithDomain(http.MethodPut, fmt.Sprintf("/domains/%s?ip=%s", domain, ip), domain)
		resp := httptest.NewRecorder()
		handler.UpdateDomain(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, fmt.Sprintf("{\"message\":\"Updated %s to %s\"}", domain, ip), resp.Body.String())
	})

	t.Run("Auto-detects IP from the connection when ip parameter is omitted", func(t *testing.T) {
		detectedIP := "203.0.113.42"
		mockService.EXPECT().Update(gomock.Any(), domain, detectedIP).Return(nil)

		req := newRequestWithDomain(http.MethodPut, fmt.Sprintf("/domains/%s", domain), domain)
		req.RemoteAddr = detectedIP + ":54321"
		resp := httptest.NewRecorder()
		handler.UpdateDomain(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, fmt.Sprintf("{\"message\":\"Updated %s to %s\"}", domain, detectedIP), resp.Body.String())
	})

	t.Run("Failed to determine IP when RemoteAddr has no port", func(t *testing.T) {
		req := newRequestWithDomain(http.MethodPut, fmt.Sprintf("/domains/%s", domain), domain)
		req.RemoteAddr = "not-a-valid-remote-addr"
		resp := httptest.NewRecorder()
		handler.UpdateDomain(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, `{"error":"Unable to determine IP"}`, resp.Body.String())
	})

	t.Run("Failed unexpected error", func(t *testing.T) {
		mockService.EXPECT().Update(gomock.Any(), domain, ip).Return(fmt.Errorf("some error"))

		req := newRequestWithDomain(http.MethodPut, fmt.Sprintf("/domains/%s?ip=%s", domain, ip), domain)
		resp := httptest.NewRecorder()
		handler.UpdateDomain(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.JSONEq(t, `{"error":"some error"}`, resp.Body.String())
	})

	t.Run("Domain not found maps to not found", func(t *testing.T) {
		mockService.EXPECT().Update(gomock.Any(), domain, ip).Return(ports.ErrDomainNotFound)

		req := newRequestWithDomain(http.MethodPut, fmt.Sprintf("/domains/%s?ip=%s", domain, ip), domain)
		resp := httptest.NewRecorder()
		handler.UpdateDomain(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		assert.JSONEq(t, `{"error":"Domain not found"}`, resp.Body.String())
	})

	t.Run("Failed validation error maps to bad request", func(t *testing.T) {
		mockService.EXPECT().Update(gomock.Any(), domain, ip).Return(ports.ErrInvalidDomain)

		req := newRequestWithDomain(http.MethodPut, fmt.Sprintf("/domains/%s?ip=%s", domain, ip), domain)
		resp := httptest.NewRecorder()
		handler.UpdateDomain(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, fmt.Sprintf(`{"error":%q}`, ports.ErrInvalidDomain.Error()), resp.Body.String())
	})
}

func TestGetDomainHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockDNSService(ctrl)
	handler := NewHandler(mockService)

	record := ports.Dns{Domain: "example.com", IP: net.ParseIP("192.168.1.1")}

	t.Run("Retrieve found DNS by domain", func(t *testing.T) {
		mockService.EXPECT().Find(gomock.Any(), record.Domain).Return(&record, nil)

		req := newRequestWithDomain(http.MethodGet, fmt.Sprintf("/domains/%s", record.Domain), record.Domain)
		resp := httptest.NewRecorder()
		handler.GetDomain(resp, req)

		expectedJSON, _ := json.Marshal(record)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, string(expectedJSON), resp.Body.String())
	})

	t.Run("Not found DNS by domain", func(t *testing.T) {
		mockService.EXPECT().Find(gomock.Any(), record.Domain).Return(nil, nil)

		req := newRequestWithDomain(http.MethodGet, fmt.Sprintf("/domains/%s", record.Domain), record.Domain)
		resp := httptest.NewRecorder()
		handler.GetDomain(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		assert.JSONEq(t, `{"error": "Domain not found"}`, resp.Body.String())
	})

	t.Run("Failed unexpected error", func(t *testing.T) {
		mockService.EXPECT().Find(gomock.Any(), record.Domain).Return(nil, fmt.Errorf("some error"))

		req := newRequestWithDomain(http.MethodGet, fmt.Sprintf("/domains/%s", record.Domain), record.Domain)
		resp := httptest.NewRecorder()
		handler.GetDomain(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.JSONEq(t, `{"error":"some error"}`, resp.Body.String())
	})
}

func TestDeleteDomainHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockDNSService(ctrl)
	handler := NewHandler(mockService)

	domain := "example.com"

	t.Run("Delete successful", func(t *testing.T) {
		mockService.EXPECT().Delete(gomock.Any(), domain).Return(nil)

		req := newRequestWithDomain(http.MethodDelete, fmt.Sprintf("/domains/%s", domain), domain)
		resp := httptest.NewRecorder()
		handler.DeleteDomain(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, fmt.Sprintf("{\"message\":\"Deleted %s\"}", domain), resp.Body.String())
	})

	t.Run("Not found DNS by domain", func(t *testing.T) {
		mockService.EXPECT().Delete(gomock.Any(), domain).Return(ports.ErrDomainNotFound)

		req := newRequestWithDomain(http.MethodDelete, fmt.Sprintf("/domains/%s", domain), domain)
		resp := httptest.NewRecorder()
		handler.DeleteDomain(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		assert.JSONEq(t, `{"error":"Domain not found"}`, resp.Body.String())
	})

	t.Run("Failed validation error maps to bad request", func(t *testing.T) {
		mockService.EXPECT().Delete(gomock.Any(), domain).Return(ports.ErrInvalidDomain)

		req := newRequestWithDomain(http.MethodDelete, fmt.Sprintf("/domains/%s", domain), domain)
		resp := httptest.NewRecorder()
		handler.DeleteDomain(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, fmt.Sprintf(`{"error":%q}`, ports.ErrInvalidDomain.Error()), resp.Body.String())
	})

	t.Run("Failed unexpected error", func(t *testing.T) {
		mockService.EXPECT().Delete(gomock.Any(), domain).Return(fmt.Errorf("some error"))

		req := newRequestWithDomain(http.MethodDelete, fmt.Sprintf("/domains/%s", domain), domain)
		resp := httptest.NewRecorder()
		handler.DeleteDomain(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.JSONEq(t, `{"error":"some error"}`, resp.Body.String())
	})
}
