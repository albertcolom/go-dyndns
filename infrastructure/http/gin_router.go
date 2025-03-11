package http

import "github.com/gin-gonic/gin"

type DnsRouter struct {
	DnsRouter *DNSHandler
}

func NewRouter(dnsHandler *DNSHandler) *gin.Engine {
	router := gin.New()

	router.GET("/update", dnsHandler.UpdateIp)
	router.GET("/get", dnsHandler.GetIp)

	return router
}
