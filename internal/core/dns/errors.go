package dns

import "fmt"

var (
	ErrDomainEmpty      = fmt.Errorf("domain cannot be empty")
	ErrInvalidDomain    = fmt.Errorf("invalid domain")
	ErrInvalidDomainLen = fmt.Errorf("domain too long (max 255 characters)")
	ErrEmptyIP          = fmt.Errorf("IP cannot be empty")
	ErrInvalidIP        = fmt.Errorf("invalid IP address")
)
