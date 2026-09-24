package service

import (
	"context"
	"net"
	"testing"

	"go-dyndns/internal/ports"
	"go-dyndns/internal/ports/mocks"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUpdateDns(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockDNSRepository(ctrl)
	service := NewDNSService(mockRepository)
	ctx := context.Background()

	t.Run("Update successful", func(t *testing.T) {
		domain := "example.com"
		ip := net.ParseIP("192.168.1.1")

		mockRepository.EXPECT().Save(ctx, &ports.Dns{Domain: domain, IP: ip}).Return(nil)
		err := service.Update(ctx, domain, ip.String())

		assert.NoError(t, err)
	})

	t.Run("Update failed for invalid IP", func(t *testing.T) {
		domain := "example.com"
		ip := "invalid ip"

		err := service.Update(ctx, domain, ip)

		assert.Error(t, err)
		assert.Equal(t, ports.ErrInvalidIP, err)
	})

	t.Run("Update failed for invalid domain", func(t *testing.T) {
		domain := "i n v a l i d .domain"
		ip := "192.168.1.1"

		err := service.Update(ctx, domain, ip)

		assert.Error(t, err)
		assert.Equal(t, ports.ErrInvalidDomain, err)
	})
}

func TestFindDns(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockDNSRepository(ctrl)
	service := NewDNSService(mockRepository)
	ctx := context.Background()

	t.Run("Retrieve found DNS", func(t *testing.T) {
		domain := "example.com"
		expected := &ports.Dns{Domain: domain, IP: net.ParseIP("192.168.1.1")}

		mockRepository.EXPECT().Find(ctx, domain).Return(expected, nil)
		result, err := service.Find(ctx, domain)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Retrieve nil when not found DNS", func(t *testing.T) {
		domain := "example.com"

		mockRepository.EXPECT().Find(ctx, domain).Return(nil, nil)
		result, err := service.Find(ctx, domain)

		assert.NoError(t, err)
		assert.Nil(t, result)
	})
}

func TestValidateDomain(t *testing.T) {
	data := []struct {
		name     string
		domain   string
		expected error
	}{
		{
			name:     "valid domain",
			domain:   "example.com",
			expected: nil,
		},
		{
			name:     "valid subdomain",
			domain:   "sub.example.co.uk",
			expected: nil,
		},
		{
			name:     "valid domain with hyphen",
			domain:   "my-domain.org",
			expected: nil,
		},
		{
			name:     "invalid domain (no TLD)",
			domain:   "example",
			expected: ports.ErrInvalidDomain,
		},
		{
			name:     "invalid domain (trailing dot)",
			domain:   "example.",
			expected: ports.ErrInvalidDomain,
		},
		{
			name:     "invalid domain (invalid characters)",
			domain:   "example!.com",
			expected: ports.ErrInvalidDomain,
		},
		{
			name:     "invalid domain (single-character TLD)",
			domain:   "example.c",
			expected: ports.ErrInvalidDomain,
		},
		{
			name:     "invalid domain (empty string)",
			domain:   "",
			expected: ports.ErrDomainEmpty,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			err := validateDomain(d.domain)
			assert.Equal(t, d.expected, err)
		})
	}
}
