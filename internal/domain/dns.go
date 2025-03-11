package domain

import (
	"fmt"
	"net"
)

type DNSRecord struct {
	Domain string `json:"domain"`
	IP     net.IP `json:"ip"`
}

func (d *DNSRecord) Validate() error {
	if d.Domain == "" {
		return ErrInvalidDomain
	}
	if d.IP == nil {
		return ErrInvalidIP
	}
	return nil
}

var (
	ErrInvalidDomain = fmt.Errorf("invalid domain")
	ErrInvalidIP     = fmt.Errorf("invalid IP address")
)
