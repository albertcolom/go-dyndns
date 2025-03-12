package dns

import (
	"fmt"
	"net"
	"regexp"
)

var (
	ErrInvalidDomain = fmt.Errorf("invalid domain")
	ErrInvalidIP     = fmt.Errorf("invalid IP address")
)

type Dns struct {
	Domain string `json:"domain"`
	IP     net.IP `json:"ip"`
}

func (d Dns) ValidateDomain() error {
	if !isValidDomain(d.Domain) {
		return ErrInvalidDomain
	}

	return nil
}

func isValidDomain(domain string) bool {
	regex := `^([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(regex, domain)

	return match
}
