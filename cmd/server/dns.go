package main

import (
	"fmt"
	"go-dyndns/internal/adapters/dns"
	"log"
)

func StartDNSServer(s *dns.Server) chan error {
	errChan := make(chan error, 1)

	go func() {
		log.Printf("[DNS] Starting server on %s (%s)", s.DnsServer.Addr, s.DnsServer.Net)
		if err := s.Start(); err != nil {
			errChan <- fmt.Errorf("[DNS] Server error: %w", err)
		}
	}()

	return errChan
}
