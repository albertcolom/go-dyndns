package main

import (
	"fmt"
	"go-dyndns/internal/adapters/http"
	"log"

	"go-dyndns/internal/core/dns"
)

func StartHTTPServer(service dns.Service, addr string) (*http.Server, chan error) {
	errChan := make(chan error, 1)
	handler := http.NewDNSHandler(service)
	httpServer := http.NewHTTPServer(handler, addr)

	go func() {
		log.Printf("Starting HTTP server on %s", addr)
		if err := httpServer.Start(); err != nil {
			errChan <- fmt.Errorf("HTTP server error: %w", err)
		}
	}()

	return httpServer, errChan
}
