//go:generate go tool mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

package port

import (
	"context"
	"fmt"
	"net"
)

const DomainPattern = `^([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`

var (
	ErrDomainEmpty      = fmt.Errorf("domain cannot be empty")
	ErrInvalidDomain    = fmt.Errorf("invalid domain")
	ErrInvalidDomainLen = fmt.Errorf("domain too long (max 255 characters)")
	ErrEmptyIP          = fmt.Errorf("IP cannot be empty")
	ErrInvalidIP        = fmt.Errorf("invalid IP address")
)

type Dns struct {
	Domain string `json:"domain"`
	IP     net.IP `json:"ip"`
}

type DNSRepository interface {
	Save(ctx context.Context, dns *Dns) error
	Find(ctx context.Context, domain string) (*Dns, error)
}

type DNSService interface {
	Update(ctx context.Context, domain, ip string) error
	Find(ctx context.Context, domain string) (*Dns, error)
}
