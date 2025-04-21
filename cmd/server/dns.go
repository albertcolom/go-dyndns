package main

import (
	"context"
	"fmt"
	"go-dyndns/internal/adapters/server"
	"go-dyndns/internal/core/dns"
	"log"
)

func StartDNSServer(ctx context.Context, service dns.Service, addr, net string) chan error {
	errChan := make(chan error, 1)
	dnsServer := server.NewDns(service)

	go func() {
		log.Printf("Starting DNS server on %s (%s)", addr, net)
		if err := dnsServer.Start(ctx, addr, net); err != nil {
			errChan <- fmt.Errorf("DNS server error: %w", err)
		}
	}()

	return errChan
}
