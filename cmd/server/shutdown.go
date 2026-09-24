package main

import (
	"context"
	"go-dyndns/internal/adapters/dns"
	"go-dyndns/internal/adapters/http"
	"go-dyndns/internal/ports"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func WaitForShutdown(
	ctx context.Context,
	cancel context.CancelFunc,
	dnsServer *dns.Server,
	httpServer *http.Server,
	httpErrChan, dnsErrChan chan error,
	log ports.Logger,
) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)

	select {
	case <-interruptChan:
		log.Info(ctx, "Received shutdown signal", "component", "SYSTEM")
	case err := <-httpErrChan:
		log.Error(ctx, "HTTP Server error", "component", "HTTP", "error", err)
	case err := <-dnsErrChan:
		log.Error(ctx, "DNS Server error", "component", "DNS", "error", err)
	}

	log.Info(ctx, "Initiating graceful shutdown", "component", "APP")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := dnsServer.Shutdown(shutdownCtx); err != nil {
			log.Error(ctx, "Server shutdown error", "component", "DNS", "error", err)
		} else {
			log.Info(ctx, "Server stopped gracefully", "component", "DNS")
		}
	}()

	go func() {
		defer wg.Done()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Error(ctx, "Server shutdown error", "component", "HTTP", "error", err)
		} else {
			log.Info(ctx, "Server stopped gracefully", "component", "HTTP")
		}
	}()

	wg.Wait()
	cancel()
}
