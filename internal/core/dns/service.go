//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE

package dns

import (
	"context"
	"net"
)

type Service interface {
	Update(ctx context.Context, domain, ip string) error
	Find(ctx context.Context, domain string) (*Dns, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (s *service) Update(ctx context.Context, domain, ip string) error {
	parseIP := net.ParseIP(ip)
	if parseIP == nil {
		return ErrInvalidIP
	}
	dns := &Dns{Domain: domain, IP: parseIP}

	if err := dns.ValidateDomain(); err != nil {
		return err
	}

	return s.repository.Save(ctx, dns)
}

func (s *service) Find(ctx context.Context, domain string) (*Dns, error) {
	return s.repository.Find(ctx, domain)
}
