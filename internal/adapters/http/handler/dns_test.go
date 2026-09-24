package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go-dyndns/internal/core/dns"
	"go.uber.org/mock/gomock"
)

func TestUpdateHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()
	defer ctrl.Finish()

	mockService := dns.NewMockService(ctrl)
	handler := NewHandler(mockService)

	domain := "example.com"
	ip := "192.168.1.1"

	t.Run("Update successful", func(t *testing.T) {
		mockService.EXPECT().Update(ctx, domain, ip).Return(nil)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/update?domain=%s&ip=%s", domain, ip), nil)
		resp := httptest.NewRecorder()
		handler.UpdateIp(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, fmt.Sprintf("{\"message\":\"Updated %s to %s\"}", domain, ip), resp.Body.String())
	})

	t.Run("Failed missing IP parameter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/update?domain=%s", domain), nil)
		resp := httptest.NewRecorder()
		handler.UpdateIp(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, `{"error":"Missing parameters"}`, resp.Body.String())
	})

	t.Run("Failed missing domain parameter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/update?ip=%s", ip), nil)
		resp := httptest.NewRecorder()
		handler.UpdateIp(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, `{"error":"Missing parameters"}`, resp.Body.String())
	})

	t.Run("Failed unexpected error", func(t *testing.T) {
		mockService.EXPECT().Update(ctx, domain, ip).Return(fmt.Errorf("some error"))

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/update?domain=%s&ip=%s", domain, ip), nil)
		resp := httptest.NewRecorder()
		handler.UpdateIp(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.JSONEq(t, `{"error":"some error"}`, resp.Body.String())
	})
}

func TestGetIpHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()
	defer ctrl.Finish()

	mockService := dns.NewMockService(ctrl)
	handler := NewHandler(mockService)

	record := dns.Dns{Domain: "example.com", IP: net.ParseIP("192.168.1.1")}

	t.Run("Retrieve found DNS by domain", func(t *testing.T) {
		mockService.EXPECT().Find(ctx, record.Domain).Return(&record, nil)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/get?domain=%s", record.Domain), nil)
		resp := httptest.NewRecorder()
		handler.GetIp(resp, req)

		expectedJSON, _ := json.Marshal(record)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, string(expectedJSON), resp.Body.String())
	})

	t.Run("Not found DNS by domain", func(t *testing.T) {
		mockService.EXPECT().Find(ctx, record.Domain).Return(nil, nil)

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
		mockService.EXPECT().Find(ctx, record.Domain).Return(nil, fmt.Errorf("some error"))

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/get?domain=%s", record.Domain), nil)
		resp := httptest.NewRecorder()
		handler.GetIp(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.JSONEq(t, `{"error":"some error"}`, resp.Body.String())
	})
}
