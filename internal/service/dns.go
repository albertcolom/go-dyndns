package service

import (
	"context"
	"net"
	"regexp"
	"time"

	"go-dyndns/internal/ports"
)

var domainPattern = regexp.MustCompile(ports.DomainPattern)

type dnsService struct {
	repository ports.DNSRepository
}

func NewDNSService(repository ports.DNSRepository) ports.DNSService {
	return &dnsService{repository: repository}
}

func (s *dnsService) Create(ctx context.Context, domain, ip string) error {
	dns := &ports.Dns{Domain: domain, IP: net.ParseIP(ip)}
	if err := validateDns(dns); err != nil {
		return err
	}

	existing, err := s.repository.Find(ctx, domain)
	if err != nil {
		return err
	}
	if existing != nil {
		return ports.ErrDomainExists
	}

	now := time.Now().UTC()
	dns.CreatedAt = now
	dns.UpdatedAt = now

	return s.repository.Save(ctx, dns)
}

func (s *dnsService) Update(ctx context.Context, domain, ip string) error {
	dns := &ports.Dns{Domain: domain, IP: net.ParseIP(ip)}
	if err := validateDns(dns); err != nil {
		return err
	}

	existing, err := s.repository.Find(ctx, domain)
	if err != nil {
		return err
	}
	if existing == nil {
		return ports.ErrDomainNotFound
	}

	dns.CreatedAt = existing.CreatedAt
	dns.UpdatedAt = time.Now().UTC()

	return s.repository.Save(ctx, dns)
}

func (s *dnsService) Find(ctx context.Context, domain string) (*ports.Dns, error) {
	return s.repository.Find(ctx, domain)
}

func (s *dnsService) Delete(ctx context.Context, domain string) error {
	if err := validateDomain(domain); err != nil {
		return err
	}
	return s.repository.Delete(ctx, domain)
}

func validateDns(d *ports.Dns) error {
	if err := validateDomain(d.Domain); err != nil {
		return err
	}
	if d.IP == nil || d.IP.To4() == nil {
		return ports.ErrInvalidIP
	}
	return nil
}

func validateDomain(domain string) error {
	if domain == "" {
		return ports.ErrDomainEmpty
	}
	if len(domain) > 255 {
		return ports.ErrInvalidDomainLen
	}
	if !domainPattern.MatchString(domain) {
		return ports.ErrInvalidDomain
	}
	return nil
}
