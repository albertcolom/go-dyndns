package http

import (
	"log"

	"github.com/gin-gonic/gin"
)

type DnsRouter struct {
	DnsHandler *DNSHandler
}

func NewRouter(dnsHandler *DNSHandler) *DnsRouter {
	return &DnsRouter{DnsHandler: dnsHandler}
}

func (r *DnsRouter) Run(addr string) {
	router := gin.New()

	router.GET("/health", r.DnsHandler.Health)
	router.GET("/update", r.DnsHandler.UpdateIp)
	router.GET("/get", r.DnsHandler.GetIp)

	log.Println("Starting HTTP server on :8080")

	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
