package main

import (
	"fmt"
	"go-dyndns/internal/adapters/http"
	"go-dyndns/pkg/logger"
)

func StartHTTPServer(s *http.Server, log logger.Logger) chan error {
	errChan := make(chan error, 1)

	go func() {
		log.Info("HTTP", "Starting server", logger.Field{Key: "addr", Value: s.HttpServer.Addr})
		if err := s.Start(); err != nil {
			errChan <- fmt.Errorf("HTTP server error: %w", err)
		}
	}()

	return errChan
}
