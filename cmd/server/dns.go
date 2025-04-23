package main

import (
	"fmt"
	server "go-dyndns/internal/adapters/dns"
	"go-dyndns/internal/core/dns"
	"log"
)

func StartDNSServer(service dns.Service, addr, net string) (*server.Server, chan error) {
	errChan := make(chan error, 1)
	dnsServer := server.NewDnsServer(service, addr, net)

	go func() {
		log.Printf("[DNS] Starting server on %s (%s)", addr, net)
		if err := dnsServer.Start(); err != nil {
			errChan <- fmt.Errorf("[DNS] Server error: %w", err)
		}
	}()

	return dnsServer, errChan
}
