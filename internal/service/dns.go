package service

import (
	"context"
	"net"
	"regexp"

	"go-dyndns/internal/ports"
)

type dnsService struct {
	repository ports.DNSRepository
}

func NewDNSService(repository ports.DNSRepository) ports.DNSService {
	return &dnsService{repository: repository}
}

func (s *dnsService) Update(ctx context.Context, domain, ip string) error {
	dns := &ports.Dns{Domain: domain, IP: net.ParseIP(ip)}
	if err := validateDns(dns); err != nil {
		return err
	}
	return s.repository.Save(ctx, dns)
}

func (s *dnsService) Find(ctx context.Context, domain string) (*ports.Dns, error) {
	return s.repository.Find(ctx, domain)
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
	match, _ := regexp.MatchString(ports.DomainPattern, domain)
	if !match {
		return ports.ErrInvalidDomain
	}
	return nil
}
