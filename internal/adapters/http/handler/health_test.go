package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go-dyndns/internal/core"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestLivezHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := core.NewMockDNSService(ctrl)
	handler := NewHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/livez", nil)
	resp := httptest.NewRecorder()
	handler.Livez(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Empty(t, resp.Body.String())
}
