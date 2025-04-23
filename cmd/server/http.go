package main

import (
	"fmt"
	"go-dyndns/internal/adapters/http"
	"log"

	"go-dyndns/internal/core/dns"
)

func StartHTTPServer(service dns.Service, addr, token string) (*http.Server, chan error) {
	errChan := make(chan error, 1)
	handler := http.NewHandler(service)
	httpServer := http.NewHTTPServer(handler, addr, token)

	go func() {
		log.Printf("[HTTP] Starting server on %s", addr)
		if err := httpServer.Start(); err != nil {
			errChan <- fmt.Errorf("[HTTP] Server error: %w", err)
		}
	}()

	return httpServer, errChan
}
