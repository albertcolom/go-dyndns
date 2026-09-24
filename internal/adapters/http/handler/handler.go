package handler

import "go-dyndns/internal/ports"

type Handler struct {
	service ports.DNSService
}

func NewHandler(service ports.DNSService) *Handler {
	return &Handler{service: service}
}
