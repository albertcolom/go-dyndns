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
	dns := &Dns{Domain: domain, IP: net.ParseIP(ip)}
	if err := dns.Validate(); err != nil {
		return err
	}
	return s.repository.Save(ctx, dns)
}

func (s *service) Find(ctx context.Context, domain string) (*Dns, error) {
	return s.repository.Find(ctx, domain)
}
