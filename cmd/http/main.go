package main

import (
	"log"

	"go-dyndns/application"
	"go-dyndns/domain"
	"go-dyndns/infrastructure/database"
	"go-dyndns/infrastructure/dns"
	"go-dyndns/infrastructure/http"
	"go-dyndns/infrastructure/repository"
)

func main() {
	dbClient, err := database.NewClient("./app.db")
	if err != nil {
		log.Fatalf("Failed to initialize database client: %v", err)
	}
	defer dbClient.Close()

	repo := repository.NewSQLiteRepository(dbClient.DB)
	domainService := domain.NewDNSService(repo)
	appService := application.NewDNSAppService(domainService)

	dnsServer := dns.NewDNSServer(appService)
	go dnsServer.Start()

	dnsHandler := http.NewDNSHandler(appService)
	router := http.NewRouter(dnsHandler)
	router.Run()
}
