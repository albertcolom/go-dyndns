package main

import (
	"context"
	"fmt"
	"go-dyndns/internal/adapters/http"
	"go-dyndns/internal/ports"
)

func StartHTTPServer(ctx context.Context, s *http.Server, log ports.Logger) chan error {
	errChan := make(chan error, 1)

	go func() {
		log.Info(ctx, "Starting server", "component", "HTTP", "addr", s.HttpServer.Addr)
		if err := s.Start(); err != nil {
			errChan <- fmt.Errorf("HTTP server error: %w", err)
		}
	}()

	return errChan
}
