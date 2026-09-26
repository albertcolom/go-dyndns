package handler

import (
	"context"
	"errors"
	"net"
	"net/http"
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

// validationErrors are the service-layer errors that stem from bad caller
// input rather than an internal failure, so they map to 400 not 500.
var validationErrors = []error{
	ports.ErrDomainEmpty,
	ports.ErrInvalidDomain,
	ports.ErrInvalidDomainLen,
	ports.ErrEmptyIP,
	ports.ErrInvalidIP,
}

func isValidationError(err error) bool {
	for _, target := range validationErrors {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

func (h *Handler) UpdateIp(w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Query().Get("domain")
	ip := r.URL.Query().Get("ip")

	if domain == "" {
		Error(w, http.StatusBadRequest, "Missing parameters")
		return
	}

	if ip == "" {
		remoteIP, err := remoteAddrIP(r)
		if err != nil {
			Error(w, http.StatusBadRequest, "Unable to determine IP")
			return
		}
		ip = remoteIP
	}

	ctx, cancel := context.WithTimeout(r.Context(), defaultTimeout)
	defer cancel()

	if err := h.service.Update(ctx, domain, ip); err != nil {
		if isValidationError(err) {
			Error(w, http.StatusBadRequest, err.Error())
			return
		}
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	JSON(w, http.StatusOK, map[string]string{"message": "Updated " + domain + " to " + ip})
}

func remoteAddrIP(r *http.Request) (string, error) {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	return host, nil
}

func (h *Handler) GetIp(w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Query().Get("domain")

	if domain == "" {
		Error(w, http.StatusBadRequest, "Missing parameters")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), defaultTimeout)
	defer cancel()

	record, err := h.service.Find(ctx, domain)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if record == nil {
		Error(w, http.StatusNotFound, "Domain not found")
		return
	}

	JSON(w, http.StatusOK, record)
}
