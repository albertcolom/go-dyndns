package main

import (
	"fmt"
	"go-dyndns/internal/adapters/dns"
	"go-dyndns/pkg/logger"
)

func StartDNSServer(s *dns.Server, log logger.Logger) chan error {
	errChan := make(chan error, 1)

	go func() {
		log.Info(
			"DNS",
			"Starting server",
			logger.Field{Key: "addr", Value: s.DnsServer.Addr},
			logger.Field{Key: "net", Value: s.DnsServer.Net},
		)
		if err := s.Start(); err != nil {
			errChan <- fmt.Errorf("DNS server error: %w", err)
		}
	}()

	return errChan
}
