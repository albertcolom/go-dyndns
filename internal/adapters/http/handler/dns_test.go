package handler

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-dyndns/internal/ports"
	"go-dyndns/internal/ports/mocks"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUpdateHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockDNSService(ctrl)
	handler := NewHandler(mockService)

	domain := "example.com"
	ip := "192.168.1.1"

	t.Run("Update successful", func(t *testing.T) {
		mockService.EXPECT().Update(gomock.Any(), domain, ip).Return(nil)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/update?domain=%s&ip=%s", domain, ip), nil)
		resp := httptest.NewRecorder()
		handler.UpdateIp(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, fmt.Sprintf("{\"message\":\"Updated %s to %s\"}", domain, ip), resp.Body.String())
	})

	t.Run("Auto-detects IP from the connection when ip parameter is omitted", func(t *testing.T) {
		detectedIP := "203.0.113.42"
		mockService.EXPECT().Update(gomock.Any(), domain, detectedIP).Return(nil)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/update?domain=%s", domain), nil)
		req.RemoteAddr = detectedIP + ":54321"
		resp := httptest.NewRecorder()
		handler.UpdateIp(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, fmt.Sprintf("{\"message\":\"Updated %s to %s\"}", domain, detectedIP), resp.Body.String())
	})

	t.Run("Failed to determine IP when RemoteAddr has no port", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/update?domain=%s", domain), nil)
		req.RemoteAddr = "not-a-valid-remote-addr"
		resp := httptest.NewRecorder()
		handler.UpdateIp(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, `{"error":"Unable to determine IP"}`, resp.Body.String())
	})

	t.Run("Failed missing domain parameter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/update?ip=%s", ip), nil)
		resp := httptest.NewRecorder()
		handler.UpdateIp(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, `{"error":"Missing parameters"}`, resp.Body.String())
	})

	t.Run("Failed unexpected error", func(t *testing.T) {
		mockService.EXPECT().Update(gomock.Any(), domain, ip).Return(fmt.Errorf("some error"))

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/update?domain=%s&ip=%s", domain, ip), nil)
		resp := httptest.NewRecorder()
		handler.UpdateIp(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.JSONEq(t, `{"error":"some error"}`, resp.Body.String())
	})

	t.Run("Failed validation error maps to bad request", func(t *testing.T) {
		mockService.EXPECT().Update(gomock.Any(), domain, ip).Return(ports.ErrInvalidDomain)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/update?domain=%s&ip=%s", domain, ip), nil)
		resp := httptest.NewRecorder()
		handler.UpdateIp(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, fmt.Sprintf(`{"error":%q}`, ports.ErrInvalidDomain.Error()), resp.Body.String())
	})
}

func TestGetIpHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockDNSService(ctrl)
	handler := NewHandler(mockService)

	record := ports.Dns{Domain: "example.com", IP: net.ParseIP("192.168.1.1")}

	t.Run("Retrieve found DNS by domain", func(t *testing.T) {
		mockService.EXPECT().Find(gomock.Any(), record.Domain).Return(&record, nil)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/get?domain=%s", record.Domain), nil)
		resp := httptest.NewRecorder()
		handler.GetIp(resp, req)

		expectedJSON, _ := json.Marshal(record)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, string(expectedJSON), resp.Body.String())
	})

	t.Run("Not found DNS by domain", func(t *testing.T) {
		mockService.EXPECT().Find(gomock.Any(), record.Domain).Return(nil, nil)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/get?domain=%s", record.Domain), nil)
		resp := httptest.NewRecorder()
		handler.GetIp(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		assert.JSONEq(t, `{"error": "Domain not found"}`, resp.Body.String())
	})

	t.Run("Failed missing domain parameter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/get", nil)
		resp := httptest.NewRecorder()
		handler.GetIp(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, `{"error":"Missing parameters"}`, resp.Body.String())
	})

	t.Run("Failed unexpected error", func(t *testing.T) {
		mockService.EXPECT().Find(gomock.Any(), record.Domain).Return(nil, fmt.Errorf("some error"))

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/get?domain=%s", record.Domain), nil)
		resp := httptest.NewRecorder()
		handler.GetIp(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.JSONEq(t, `{"error":"some error"}`, resp.Body.String())
	})
}
