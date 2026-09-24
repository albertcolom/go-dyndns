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

func (d *Dns) Validate() error {
	if err := d.validateDomain(); err != nil {
		return err
	}
	if d.IP == nil || d.IP.To4() == nil {
		return ErrInvalidIP
	}
	return nil
}

func (d *Dns) validateDomain() error {
	if d.Domain == "" {
		return ErrDomainEmpty
	}
	if len(d.Domain) > 255 {
		return ErrInvalidDomainLen
	}
	match, _ := regexp.MatchString(domainPattern, d.Domain)
	if !match {
		return ErrInvalidDomain
	}
	return nil
}
