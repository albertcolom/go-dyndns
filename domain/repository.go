package domain

type DNSRepository interface {
	Save(record DNSRecord) error
	Find(domain string) (*DNSRecord, error)
}
