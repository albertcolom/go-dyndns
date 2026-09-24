package service

import (
	"context"
	"net"
	"regexp"

	"go-dyndns/internal/port"
)

type dnsService struct {
	repository port.DNSRepository
}

func NewDNSService(repository port.DNSRepository) port.DNSService {
	return &dnsService{repository: repository}
}

func (s *dnsService) Update(ctx context.Context, domain, ip string) error {
	dns := &port.Dns{Domain: domain, IP: net.ParseIP(ip)}
	if err := validateDns(dns); err != nil {
		return err
	}
	return s.repository.Save(ctx, dns)
}

func (s *dnsService) Find(ctx context.Context, domain string) (*port.Dns, error) {
	return s.repository.Find(ctx, domain)
}

func validateDns(d *port.Dns) error {
	if err := validateDomain(d.Domain); err != nil {
		return err
	}
	if d.IP == nil || d.IP.To4() == nil {
		return port.ErrInvalidIP
	}
	return nil
}

func validateDomain(domain string) error {
	if domain == "" {
		return port.ErrDomainEmpty
	}
	if len(domain) > 255 {
		return port.ErrInvalidDomainLen
	}
	match, _ := regexp.MatchString(port.DomainPattern, domain)
	if !match {
		return port.ErrInvalidDomain
	}
	return nil
}
