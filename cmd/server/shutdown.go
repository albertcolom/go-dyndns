package main

import (
	"context"
	"go-dyndns/internal/adapters/dns"
	"go-dyndns/internal/adapters/http"
	"go-dyndns/pkg/logger"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func WaitForShutdown(
	cancel context.CancelFunc,
	dnsServer *dns.Server,
	httpServer *http.Server,
	httpErrChan, dnsErrChan chan error,
	log logger.Logger,
) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)

	select {
	case <-interruptChan:
		log.Info("SYSTEM", "Received shutdown signal")
	case err := <-httpErrChan:
		log.Error("HTTP", "HTTP Server error", logger.Field{Key: "error", Value: err})
	case err := <-dnsErrChan:
		log.Error("DNS", "DNS Server error", logger.Field{Key: "error", Value: err})
	}

	log.Info("APP", "Initiating graceful shutdown")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := dnsServer.Shutdown(shutdownCtx); err != nil {
			log.Error("DNS", "Server shutdown error", logger.Field{Key: "error", Value: err})
		} else {
			log.Info("DNS", "Server stopped gracefully")
		}
	}()

	go func() {
		defer wg.Done()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Error("HTTP", "Server shutdown error", logger.Field{Key: "error", Value: err})
		} else {
			log.Info("HTTP", "Server stopped gracefully")
		}
	}()

	wg.Wait()
	cancel()
}
