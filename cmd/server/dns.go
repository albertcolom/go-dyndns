package main

import (
	"context"
	"fmt"
	"go-dyndns/internal/adapters/dns"
	"go-dyndns/internal/port"
)

func StartDNSServer(ctx context.Context, s *dns.Server, log port.Logger) chan error {
	errChan := make(chan error, 1)

	go func() {
		log.Info(ctx, "Starting server",
			"component", "DNS",
			"addr", s.DnsServer.Addr,
			"net", s.DnsServer.Net,
		)
		if err := s.Start(); err != nil {
			errChan <- fmt.Errorf("DNS server error: %w", err)
		}
	}()

	return errChan
}
