//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE

package dns

import "net"

type Service interface {
	Update(domain, ip string) error
	Find(domain string) (*Dns, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (s *service) Update(domain, ip string) error {
	parseIP := net.ParseIP(ip)
	if parseIP == nil {
		return ErrInvalidIP
	}
	dns := &Dns{Domain: domain, IP: parseIP}

	if err := dns.ValidateDomain(); err != nil {
		return err
	}

	return s.repository.Save(dns)
}

func (s *service) Find(domain string) (*Dns, error) {
	return s.repository.Find(domain)
}
