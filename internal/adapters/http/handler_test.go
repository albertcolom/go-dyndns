package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go-dyndns/internal/core/dns"
	"go.uber.org/mock/gomock"
)

func TestHealthHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	router := gin.Default()
	defer ctrl.Finish()

	mockService := dns.NewMockService(ctrl)
	handler := NewHandler(mockService)

	router.GET("/health", handler.Health)

	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.JSONEq(t, `{"status":"ok"}`, resp.Body.String())
}

func TestUpdateHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	router := gin.Default()
	ctx := context.Background()
	defer ctrl.Finish()

	mockService := dns.NewMockService(ctrl)
	handler := NewHandler(mockService)
	router.GET("/update", handler.UpdateIp)

	domain := "example.com"
	ip := "192.168.1.1"

	t.Run("Update successful", func(t *testing.T) {
		mockService.EXPECT().Update(ctx, domain, ip).Return(nil)

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/update?domain=%s&ip=%s", domain, ip), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, fmt.Sprintf("{\"message\":\"Updated %s to %s\"}", domain, ip), resp.Body.String())
	})

	t.Run("Failed missing IP parameter", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/update?domain=%s", domain), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, `{"error":"Missing parameters"}`, resp.Body.String())
	})

	t.Run("Failed missing domain parameter", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/update?ip=%s", ip), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, `{"error":"Missing parameters"}`, resp.Body.String())
	})

	t.Run("Failed unexpected error", func(t *testing.T) {
		mockService.EXPECT().Update(ctx, domain, ip).Return(fmt.Errorf("some error"))

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/update?domain=%s&ip=%s", domain, ip), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.JSONEq(t, `{"error":"some error"}`, resp.Body.String())
	})
}

func TestGetIpHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	router := gin.Default()
	ctx := context.Background()
	defer ctrl.Finish()

	mockService := dns.NewMockService(ctrl)
	handler := NewHandler(mockService)
	router.GET("/get", handler.GetIp)

	dns := dns.Dns{Domain: "example.com", IP: net.ParseIP("192.168.1.1")}

	t.Run("Retrieve found DNS by domain", func(t *testing.T) {
		mockService.EXPECT().Find(ctx, dns.Domain).Return(&dns, nil)

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/get?domain=%s", dns.Domain), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		expectedJSON, _ := json.Marshal(dns)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, string(expectedJSON), resp.Body.String())
	})

	t.Run("Not found DNS by domain", func(t *testing.T) {
		mockService.EXPECT().Find(ctx, dns.Domain).Return(nil, nil)

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/get?domain=%s", dns.Domain), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		assert.JSONEq(t, `{"error": "Domain not found"}`, resp.Body.String())
	})

	t.Run("Failed missing domain parameter", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/get", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, `{"error":"Missing parameters"}`, resp.Body.String())
	})

	t.Run("Failed unexpected error", func(t *testing.T) {
		mockService.EXPECT().Find(ctx, dns.Domain).Return(nil, fmt.Errorf("some error"))

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/get?domain=%s", dns.Domain), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.JSONEq(t, `{"error":"some error"}`, resp.Body.String())
	})
}
