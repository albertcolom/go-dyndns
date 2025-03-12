//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE

package dns

type Repository interface {
	Save(dns *Dns) error
	Find(domain string) (*Dns, error)
}
