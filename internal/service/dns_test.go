package service

import (
	"context"
	"net"
	"testing"
	"time"

	"go-dyndns/internal/ports"
	"go-dyndns/internal/ports/mocks"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateDns(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockDNSRepository(ctrl)
	service := NewDNSService(mockRepository)
	ctx := context.Background()

	t.Run("Create successful sets CreatedAt and UpdatedAt to the same fresh timestamp", func(t *testing.T) {
		domain := "example.com"
		ip := net.ParseIP("192.168.1.1")

		mockRepository.EXPECT().Find(ctx, domain).Return(nil, nil)
		mockRepository.EXPECT().Save(ctx, gomock.Any()).DoAndReturn(func(_ context.Context, saved *ports.Dns) error {
			assert.Equal(t, domain, saved.Domain)
			assert.True(t, ip.Equal(saved.IP))
			assert.WithinDuration(t, time.Now(), saved.CreatedAt, time.Second)
			assert.Equal(t, saved.CreatedAt, saved.UpdatedAt)
			return nil
		})
		err := service.Create(ctx, domain, ip.String())

		assert.NoError(t, err)
	})

	t.Run("Create fails when domain already exists", func(t *testing.T) {
		domain := "example.com"
		ip := "192.168.1.1"
		existing := &ports.Dns{Domain: domain, IP: net.ParseIP("10.0.0.1")}

		mockRepository.EXPECT().Find(ctx, domain).Return(existing, nil)
		err := service.Create(ctx, domain, ip)

		assert.Equal(t, ports.ErrDomainExists, err)
	})

	t.Run("Create failed for invalid IP", func(t *testing.T) {
		domain := "example.com"
		ip := "invalid ip"

		err := service.Create(ctx, domain, ip)

		assert.Error(t, err)
		assert.Equal(t, ports.ErrInvalidIP, err)
	})

	t.Run("Create failed for invalid domain", func(t *testing.T) {
		domain := "i n v a l i d .domain"
		ip := "192.168.1.1"

		err := service.Create(ctx, domain, ip)

		assert.Error(t, err)
		assert.Equal(t, ports.ErrInvalidDomain, err)
	})
}

func TestUpdateDns(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockDNSRepository(ctrl)
	service := NewDNSService(mockRepository)
	ctx := context.Background()

	t.Run("Update successful preserves CreatedAt and refreshes UpdatedAt", func(t *testing.T) {
		domain := "example.com"
		ip := net.ParseIP("192.168.1.1")
		existingCreatedAt := time.Now().Add(-time.Hour).UTC()
		existing := &ports.Dns{Domain: domain, IP: net.ParseIP("10.0.0.1"), CreatedAt: existingCreatedAt, UpdatedAt: existingCreatedAt}

		mockRepository.EXPECT().Find(ctx, domain).Return(existing, nil)
		mockRepository.EXPECT().Save(ctx, gomock.Any()).DoAndReturn(func(_ context.Context, saved *ports.Dns) error {
			assert.Equal(t, domain, saved.Domain)
			assert.True(t, ip.Equal(saved.IP))
			assert.True(t, existingCreatedAt.Equal(saved.CreatedAt))
			assert.WithinDuration(t, time.Now(), saved.UpdatedAt, time.Second)
			assert.True(t, saved.UpdatedAt.After(existingCreatedAt))
			return nil
		})
		err := service.Update(ctx, domain, ip.String())

		assert.NoError(t, err)
	})

	t.Run("Update fails when domain does not exist", func(t *testing.T) {
		domain := "example.com"
		ip := "192.168.1.1"

		mockRepository.EXPECT().Find(ctx, domain).Return(nil, nil)
		err := service.Update(ctx, domain, ip)

		assert.Equal(t, ports.ErrDomainNotFound, err)
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

func TestDeleteDns(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockDNSRepository(ctrl)
	service := NewDNSService(mockRepository)
	ctx := context.Background()

	t.Run("Delete successful", func(t *testing.T) {
		domain := "example.com"

		mockRepository.EXPECT().Delete(ctx, domain).Return(nil)
		err := service.Delete(ctx, domain)

		assert.NoError(t, err)
	})

	t.Run("Delete propagates not found from repository", func(t *testing.T) {
		domain := "example.com"

		mockRepository.EXPECT().Delete(ctx, domain).Return(ports.ErrDomainNotFound)
		err := service.Delete(ctx, domain)

		assert.Equal(t, ports.ErrDomainNotFound, err)
	})

	t.Run("Delete failed for invalid domain", func(t *testing.T) {
		domain := "i n v a l i d .domain"

		err := service.Delete(ctx, domain)

		assert.Error(t, err)
		assert.Equal(t, ports.ErrInvalidDomain, err)
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
