package main

import (
	"log"

	"go-dyndns/internal/application"
	"go-dyndns/internal/domain"
	"go-dyndns/internal/infrastructure/config"
	"go-dyndns/internal/infrastructure/database"
	"go-dyndns/internal/infrastructure/dns"
	"go-dyndns/internal/infrastructure/http"
	"go-dyndns/internal/infrastructure/repository"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dbClient, err := database.NewClient(cfg.Sqlite.Path)
	if err != nil {
		log.Fatalf("Failed to initialize database client: %v", err)
	}
	defer dbClient.Close()

	repo := repository.NewSQLiteRepository(dbClient.DB)
	domainService := domain.NewDNSService(repo)
	appService := application.NewDNSAppService(domainService)

	dnsServer := dns.NewDNSServer(appService)
	go dnsServer.Start(cfg.Dns.Addr, cfg.Dns.Net)

	dnsHandler := http.NewDNSHandler(appService)
	router := http.NewRouter(dnsHandler)
	router.Run(cfg.Http.Addr)
}
