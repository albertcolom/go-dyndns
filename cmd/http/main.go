package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"go-dyndns/application"
	"go-dyndns/domain"
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

	dnsHandler := http.NewDNSHandler(appService)

	router := gin.Default()

	router.GET("/update", dnsHandler.UpdateIp)
	router.GET("/get", dnsHandler.GetIp)

	log.Println("Starting HTTP server on :8080")
	router.Run(":8080")
}
