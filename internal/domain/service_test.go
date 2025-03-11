package domain

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUpdateDNSRecord(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := NewMockDNSRepository(ctrl)
	service := NewDNSService(mockRepository)

	t.Run("Success", func(t *testing.T) {
		domain := "example.com"
		ip := net.ParseIP("192.168.1.1")

		mockRepository.EXPECT().Save(DNSRecord{Domain: domain, IP: ip}).Return(nil)
		err := service.UpdateDNSRecord(domain, ip)

		assert.NoError(t, err)
	})

	t.Run("Invalid Domain", func(t *testing.T) {
		domain := ""
		ip := net.ParseIP("192.168.1.1")

		err := service.UpdateDNSRecord(domain, ip)

		assert.Error(t, err)
		assert.Equal(t, ErrInvalidDomain, err)
	})

	t.Run("Invalid IP", func(t *testing.T) {
		domain := "example.com"
		ip := net.ParseIP("")

		err := service.UpdateDNSRecord(domain, ip)

		assert.Error(t, err)
		assert.Equal(t, ErrInvalidIP, err)
	})
}

func TestDeleteDNSRecord(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := NewMockDNSRepository(ctrl)
	service := NewDNSService(mockRepository)

	t.Run("Success", func(t *testing.T) {
		domain := "example.com"
		expected := &DNSRecord{Domain: domain, IP: net.ParseIP("192.168.1.1")}

		mockRepository.EXPECT().Find(domain).Return(expected, nil)

		result, err := service.GetDNSRecord(domain)

		assert.Nil(t, err)
		assert.Equal(t, expected, result)
	})
}
