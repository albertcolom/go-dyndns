package handler

import "go-dyndns/internal/port"

type Handler struct {
	service port.DNSService
}

func NewHandler(service port.DNSService) *Handler {
	return &Handler{service: service}
}
