package http

import (
	"log"

	"github.com/gin-gonic/gin"
)

type Router struct {
	DnsHandler *DNSHandler
}

func NewRouter(dnsHandler *DNSHandler) *Router {
	return &Router{DnsHandler: dnsHandler}
}

func (r *Router) Run(addr string) {
	router := gin.New()

	router.GET("/health", r.DnsHandler.Health)
	router.GET("/update", r.DnsHandler.UpdateIp)
	router.GET("/get", r.DnsHandler.GetIp)

	log.Println("Starting HTTP server on :8080")

	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
