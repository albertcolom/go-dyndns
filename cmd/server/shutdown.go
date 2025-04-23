package main

import (
	"context"
	"go-dyndns/internal/adapters/dns"
	"go-dyndns/internal/adapters/http"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func WaitForShutdown(cancel context.CancelFunc, dnsServer *dns.Server, httpServer *http.Server, httpErrChan, dnsErrChan chan error) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)

	select {
	case <-interruptChan:
		log.Println("Received shutdown signal, initiating graceful shutdown...")
	case err := <-httpErrChan:
		log.Printf("Server error: %v, initiating shutdown...", err)
	case err := <-dnsErrChan:
		log.Printf("Server error: %v, initiating shutdown...", err)
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := dnsServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("[DNS] Server shutdown error: %v", err)
		} else {
			log.Println("[DNS] Server stopped gracefully")
		}
	}()

	go func() {
		defer wg.Done()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("[HTTP] Server shutdown error: %v", err)
		} else {
			log.Println("[HTTP] Server stopped gracefully")
		}
	}()

	wg.Wait()
	cancel()
}
