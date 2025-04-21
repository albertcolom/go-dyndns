package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-dyndns/internal/core/dns"
)

type Handler struct {
	service dns.Service
}

func NewHandler(service dns.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) UpdateIp(c *gin.Context) {
	domain := c.Query("domain")
	ip := c.Query("ip")

	if domain == "" || ip == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing parameters"})
		return
	}

	err := h.service.Update(c.Request.Context(), domain, ip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated " + domain + " to " + ip})
}

func (h *Handler) GetIp(c *gin.Context) {
	domain := c.Query("domain")

	if domain == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing parameters"})
		return
	}

	record, err := h.service.Find(c.Request.Context(), domain)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if record == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Domain not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"domain": record.Domain, "ip": record.IP.String()})
}
