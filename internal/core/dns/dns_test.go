package dns

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateDomain(t *testing.T) {
	data := []struct {
		name     string
		domain   string
		expected error
	}{
		{
			name:     "valid domain",
			domain:   "example.com",
			expected: nil,
		},
		{
			name:     "valid subdomain",
			domain:   "sub.example.co.uk",
			expected: nil,
		},
		{
			name:     "valid domain with hyphen",
			domain:   "my-domain.org",
			expected: nil,
		},
		{
			name:     "invalid domain (no TLD)",
			domain:   "example",
			expected: ErrInvalidDomain,
		},
		{
			name:     "invalid domain (trailing dot)",
			domain:   "example.",
			expected: ErrInvalidDomain,
		},
		{
			name:     "invalid domain (invalid characters)",
			domain:   "example!.com",
			expected: ErrInvalidDomain,
		},
		{
			name:     "invalid domain (single-character TLD)",
			domain:   "example.c",
			expected: ErrInvalidDomain,
		},
		{
			name:     "invalid domain (empty string)",
			domain:   "",
			expected: ErrDomainEmpty,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			dns := Dns{Domain: d.domain, IP: net.ParseIP("127.0.0.1")}
			err := dns.Validate()
			assert.Equal(t, d.expected, err)
		})
	}
}
