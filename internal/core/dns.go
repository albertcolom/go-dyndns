//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE

package core

import (
	"context"
	"fmt"
	"net"
)

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
