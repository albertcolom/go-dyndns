package core

import (
	"context"
	"net"
	"regexp"
)

const domainPattern = `^([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`

type service struct {
	repository DNSRepository
}

func NewService(repository DNSRepository) DNSService {
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

func (d *Dns) Validate() error {
	if err := d.validateDomain(); err != nil {
		return err
	}
	if d.IP == nil || d.IP.To4() == nil {
		return ErrInvalidIP
	}
	return nil
}

func (d *Dns) validateDomain() error {
	if d.Domain == "" {
		return ErrDomainEmpty
	}
	if len(d.Domain) > 255 {
		return ErrInvalidDomainLen
	}
	match, _ := regexp.MatchString(domainPattern, d.Domain)
	if !match {
		return ErrInvalidDomain
	}
	return nil
}
