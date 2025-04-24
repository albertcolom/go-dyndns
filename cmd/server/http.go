package main

import (
	"fmt"
	"go-dyndns/internal/adapters/http"
	"log"
)

func StartHTTPServer(s *http.Server) chan error {
	errChan := make(chan error, 1)

	go func() {
		log.Printf("[HTTP] Starting server on %s", s.HttpServer.Addr)
		if err := s.Start(); err != nil {
			errChan <- fmt.Errorf("[HTTP] Server error: %w", err)
		}
	}()

	return errChan
}
