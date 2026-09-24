package handler

import "go-dyndns/internal/core"

type Handler struct {
	service dns.Service
}

func NewHandler(service dns.Service) *Handler {
	return &Handler{service: service}
}
