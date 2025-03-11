package domain

import "net"

type DNSService struct {
	repository DNSRepository
}

func NewDNSService(repository DNSRepository) *DNSService {
	return &DNSService{repository: repository}
}

func (s *DNSService) UpdateDNSRecord(domain string, ip net.IP) error {
	record := DNSRecord{
		Domain: domain,
		IP:     ip,
	}
	if err := record.Validate(); err != nil {
		return err
	}
	return s.repository.Save(record)
}

func (s *DNSService) GetDNSRecord(domain string) (*DNSRecord, error) {
	return s.repository.Find(domain)
}
