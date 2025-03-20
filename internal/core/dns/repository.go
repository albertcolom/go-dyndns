//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE

package dns

import "context"

type Repository interface {
	Save(ctx context.Context, dns *Dns) error
	Find(ctx context.Context, domain string) (*Dns, error)
}
