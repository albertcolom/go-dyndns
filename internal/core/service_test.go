package dns

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUpdateDns(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := NewMockRepository(ctrl)
	service := NewService(mockRepository)
	ctx := context.Background()

	t.Run("Update successful", func(t *testing.T) {
		domain := "example.com"
		ip := net.ParseIP("192.168.1.1")

		mockRepository.EXPECT().Save(ctx, &Dns{Domain: domain, IP: ip}).Return(nil)
		err := service.Update(ctx, domain, ip.String())

		assert.NoError(t, err)
	})

	t.Run("Update failed for invalid IP", func(t *testing.T) {
		domain := "example.com"
		ip := "invalid ip"

		err := service.Update(ctx, domain, ip)

		assert.Error(t, err)
		assert.Equal(t, ErrInvalidIP, err)
	})

	t.Run("Update failed for invalid domain", func(t *testing.T) {
		domain := "i n v a l i d .domain"
		ip := "192.168.1.1"

		err := service.Update(ctx, domain, ip)

		assert.Error(t, err)
		assert.Equal(t, ErrInvalidDomain, err)
	})
}

func TestFindDns(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := NewMockRepository(ctrl)
	service := NewService(mockRepository)
	ctx := context.Background()

	t.Run("Retrieve found DNS", func(t *testing.T) {
		domain := "example.com"
		expected := &Dns{Domain: domain, IP: net.ParseIP("192.168.1.1")}

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
