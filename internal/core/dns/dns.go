package dns

import (
	"net"
	"regexp"
)

const domainPattern = `^([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`

type Dns struct {
	Domain string `json:"domain"`
	IP     net.IP `json:"ip"`
}

func (d *Dns) ValidateDomain() error {
	match, _ := regexp.MatchString(domainPattern, d.Domain)
	if !match {
		return ErrInvalidDomain
	}
	return nil
}
