package main

import (
	"log"

	"go-dyndns/internal/adapters/config"
	"go-dyndns/internal/adapters/http"
	"go-dyndns/internal/adapters/repository"
	"go-dyndns/internal/adapters/server"
	"go-dyndns/internal/core/dns"
	"go-dyndns/pkg/db"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dbClient, err := db.NewSQLiteClient(cfg.Sqlite.Path)
	if err != nil {
		log.Fatalf("Failed to initialize database client: %v", err)
	}
	defer dbClient.Close()

	repo := repository.NewSQLiteDNSRepository(dbClient.DB)
	service := dns.NewService(repo)

	dnsServer := server.NewDns(service)
	go dnsServer.Start(cfg.Dns.Addr, cfg.Dns.Net)

	handler := http.NewDNSHandler(service)
	router := http.NewRouter(handler)

	router.Run(cfg.Http.Addr)
}
