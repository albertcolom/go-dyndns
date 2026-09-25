package handler

import (
	"time"

	"go-dyndns/internal/ports"
)

const defaultTimeout = 2 * time.Second

type Handler struct {
	service ports.DNSService
}

func NewHandler(service ports.DNSService) *Handler {
	return &Handler{service: service}
}
