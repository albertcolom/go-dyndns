package main

import (
	"context"
	"go-dyndns/internal/adapters/config"
	server "go-dyndns/internal/adapters/dns"
	"go-dyndns/internal/adapters/http"
	"go-dyndns/internal/adapters/repository"
	"go-dyndns/internal/core/dns"
	"go-dyndns/pkg/db"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbClient, err := db.NewSQLiteClient(cfg.Sqlite.Path)
	if err != nil {
		log.Fatalf("Failed to initialize database client: %v", err)
	}
	defer func() {
		if err := dbClient.Close(); err != nil {
			log.Printf("Database close error: %v", err)
		}
	}()

	repo := repository.NewSQLiteDNSRepository(dbClient.DB)
	service := dns.NewService(repo)

	dnsServer := server.NewDnsServer(service, cfg.Dns.Addr, cfg.Dns.Net)
	dnsErrChan := StartDNSServer(dnsServer)

	httpHandler := http.NewHandler(service)
	httpServer := http.NewHTTPServer(httpHandler, cfg.Http.Addr, cfg.Http.Token)
	httpErrChan := StartHTTPServer(httpServer)

	WaitForShutdown(cancel, dnsServer, httpServer, httpErrChan, dnsErrChan)
}
