package handler

import "go-dyndns/internal/core"

type Handler struct {
	service core.DNSService
}

func NewHandler(service core.DNSService) *Handler {
	return &Handler{service: service}
}
