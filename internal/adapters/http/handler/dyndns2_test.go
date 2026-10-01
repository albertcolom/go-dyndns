package handler

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"go-dyndns/internal/ports"
	"go-dyndns/internal/ports/mocks"
)

func TestDynDNS2UpdateHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockDNSService(ctrl)
	handler := NewHandler(mockService)

	domain := "home.example.com"
	ip := "192.168.1.1"

	t.Run("Creates a new hostname and reports good", func(t *testing.T) {
		mockService.EXPECT().Find(gomock.Any(), domain).Return(nil, nil)
		mockService.EXPECT().Create(gomock.Any(), domain, ip).Return(nil)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/nic/update?hostname=%s&myip=%s", domain, ip), nil)
		resp := httptest.NewRecorder()
		handler.DynDNS2Update(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, "good "+ip+"\n", resp.Body.String())
	})

	t.Run("Updates an existing hostname with a changed IP and reports good", func(t *testing.T) {
		mockService.EXPECT().Find(gomock.Any(), domain).Return(&ports.Dns{Domain: domain, IP: net.ParseIP("10.0.0.1")}, nil)
		mockService.EXPECT().Update(gomock.Any(), domain, ip).Return(nil)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/nic/update?hostname=%s&myip=%s", domain, ip), nil)
		resp := httptest.NewRecorder()
		handler.DynDNS2Update(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, "good "+ip+"\n", resp.Body.String())
	})

	t.Run("Unchanged IP reports nochg without writing", func(t *testing.T) {
		mockService.EXPECT().Find(gomock.Any(), domain).Return(&ports.Dns{Domain: domain, IP: net.ParseIP(ip)}, nil)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/nic/update?hostname=%s&myip=%s", domain, ip), nil)
		resp := httptest.NewRecorder()
		handler.DynDNS2Update(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, "nochg "+ip+"\n", resp.Body.String())
	})

	t.Run("Auto-detects IP from the connection when myip is omitted", func(t *testing.T) {
		detectedIP := "203.0.113.42"
		mockService.EXPECT().Find(gomock.Any(), domain).Return(nil, nil)
		mockService.EXPECT().Create(gomock.Any(), domain, detectedIP).Return(nil)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/nic/update?hostname=%s", domain), nil)
		req.RemoteAddr = detectedIP + ":54321"
		resp := httptest.NewRecorder()
		handler.DynDNS2Update(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, "good "+detectedIP+"\n", resp.Body.String())
	})

	t.Run("Missing hostname reports notfqdn", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/nic/update?myip="+ip, nil)
		resp := httptest.NewRecorder()
		handler.DynDNS2Update(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, "notfqdn\n", resp.Body.String())
	})

	t.Run("Undeterminable IP reports 911", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/nic/update?hostname="+domain, nil)
		req.RemoteAddr = "not-a-valid-remote-addr"
		resp := httptest.NewRecorder()
		handler.DynDNS2Update(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, "911\n", resp.Body.String())
	})

	t.Run("Invalid IP reports 911", func(t *testing.T) {
		mockService.EXPECT().Find(gomock.Any(), domain).Return(nil, nil)
		mockService.EXPECT().Create(gomock.Any(), domain, "not-an-ip").Return(ports.ErrInvalidIP)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/nic/update?hostname=%s&myip=not-an-ip", domain), nil)
		resp := httptest.NewRecorder()
		handler.DynDNS2Update(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, "911\n", resp.Body.String())
	})

	t.Run("Invalid hostname reports notfqdn", func(t *testing.T) {
		badDomain := "i n v a l i d"
		mockService.EXPECT().Find(gomock.Any(), badDomain).Return(nil, nil)
		mockService.EXPECT().Create(gomock.Any(), badDomain, ip).Return(ports.ErrInvalidDomain)

		req := httptest.NewRequest(http.MethodGet, "/nic/update?hostname="+url.QueryEscape(badDomain)+"&myip="+ip, nil)
		resp := httptest.NewRecorder()
		handler.DynDNS2Update(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, "notfqdn\n", resp.Body.String())
	})

	t.Run("Unexpected error reports 911", func(t *testing.T) {
		mockService.EXPECT().Find(gomock.Any(), domain).Return(nil, fmt.Errorf("some error"))

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/nic/update?hostname=%s&myip=%s", domain, ip), nil)
		resp := httptest.NewRecorder()
		handler.DynDNS2Update(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Equal(t, "911\n", resp.Body.String())
	})
}
