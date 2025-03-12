package dns

import "fmt"

var (
	ErrInvalidDomain = fmt.Errorf("invalid domain")
	ErrInvalidIP     = fmt.Errorf("invalid IP address")
)
