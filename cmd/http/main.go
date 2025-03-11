package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"go-dyndns/application"
	"go-dyndns/domain"
	"go-dyndns/infrastructure/dns"
	"go-dyndns/infrastructure/http"
	"go-dyndns/infrastructure/repository"
)

func main() {
	db, err := sql.Open("sqlite3", "./app.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	repo := repository.NewSQLiteRepository(db)
	domainService := domain.NewDNSService(repo)
	appService := application.NewDNSAppService(domainService)

	dnsServer := dns.NewDNSServer(repo)
	go dnsServer.Start()

	dnsHandler := http.NewDNSHandler(appService)
	router := http.NewRouter(dnsHandler)

	log.Println("Starting HTTP server on :8080")

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
