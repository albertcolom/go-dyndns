package dns

import "net"

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Update(domain, ip string) error {
	parseIP := net.ParseIP(ip)
	if parseIP == nil {
		return ErrInvalidIP
	}
	dns := &Dns{Domain: domain, IP: parseIP}

	return s.repository.Save(dns)
}

func (s *Service) Find(domain string) (*Dns, error) {
	return s.repository.Find(domain)
}
