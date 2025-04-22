package main

import (
	"context"
	"go-dyndns/internal/adapters/config"
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

	dnsServer, dnsErrChan := StartDNSServer(service, cfg.Dns.Addr, cfg.Dns.Net)
	httpServer, httpErrChan := StartHTTPServer(service, cfg.Http.Addr, cfg.Http.Token)

	WaitForShutdown(cancel, dnsServer, httpServer, httpErrChan, dnsErrChan)
}
