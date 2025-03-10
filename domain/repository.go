//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE

package domain

type DNSRepository interface {
	Save(record DNSRecord) error
	Find(domainName string) (*DNSRecord, error)
}
