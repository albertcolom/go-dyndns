package application

import (
	"net"

	"go-dyndns/internal/domain"
)

type DNSAppService struct {
	domainService *domain.DNSService
}

func NewDNSAppService(domainService *domain.DNSService) *DNSAppService {
	return &DNSAppService{domainService}
}

func (s *DNSAppService) UpdateDNSRecord(domainName, ip string) error {
	parseIP := net.ParseIP(ip)
	if parseIP == nil {
		return domain.ErrInvalidDomain
	}
	return s.domainService.UpdateDNSRecord(domainName, parseIP)
}

func (s *DNSAppService) GetDNSRecord(domainName string) (*domain.DNSRecord, error) {
	return s.domainService.GetDNSRecord(domainName)
}
