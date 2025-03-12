package dns

import (
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

	t.Run("Update successful", func(t *testing.T) {
		domain := "example.com"
		ip := net.ParseIP("192.168.1.1")

		mockRepository.EXPECT().Save(&Dns{Domain: domain, IP: ip}).Return(nil)
		err := service.Update(domain, ip.String())

		assert.NoError(t, err)
	})

	t.Run("Update failed for invalid IP", func(t *testing.T) {
		domain := "example.com"
		ip := "invalid ip"

		err := service.Update(domain, ip)

		assert.Error(t, err)
		assert.Equal(t, ErrInvalidIP, err)
	})
}

func TestFindDns(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := NewMockRepository(ctrl)
	service := NewService(mockRepository)

	t.Run("Retrieve found DNS", func(t *testing.T) {
		domain := "example.com"
		expected := &Dns{Domain: domain, IP: net.ParseIP("192.168.1.1")}

		mockRepository.EXPECT().Find(domain).Return(expected, nil)
		result, err := service.Find(domain)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Retrieve nil when not found DNS", func(t *testing.T) {
		domain := "example.com"

		mockRepository.EXPECT().Find(domain).Return(nil, nil)
		result, err := service.Find(domain)

		assert.NoError(t, err)
		assert.Nil(t, result)
	})
}
