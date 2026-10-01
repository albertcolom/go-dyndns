package handler

import (
	"context"
	"errors"
	"net/http"

	"go-dyndns/internal/ports"
)

func (h *Handler) DynDNS2Update(w http.ResponseWriter, r *http.Request) {
	hostname := r.URL.Query().Get("hostname")
	if hostname == "" {
		Text(w, http.StatusBadRequest, "notfqdn")
		return
	}

	ip := r.URL.Query().Get("myip")
	if ip == "" {
		remoteIP, err := remoteAddrIP(r)
		if err != nil {
			Text(w, http.StatusBadRequest, "911")
			return
		}
		ip = remoteIP
	}

	ctx, cancel := context.WithTimeout(r.Context(), defaultTimeout)
	defer cancel()

	existing, err := h.service.Find(ctx, hostname)
	if err != nil {
		Text(w, http.StatusInternalServerError, "911")
		return
	}

	if existing != nil && existing.IP.String() == ip {
		Text(w, http.StatusOK, "nochg "+ip)
		return
	}

	if existing == nil {
		err = h.service.Create(ctx, hostname, ip)
	} else {
		err = h.service.Update(ctx, hostname, ip)
	}

	switch {
	case err == nil:
		Text(w, http.StatusOK, "good "+ip)
	case errors.Is(err, ports.ErrInvalidIP), errors.Is(err, ports.ErrEmptyIP):
		Text(w, http.StatusBadRequest, "911")
	case isValidationError(err):
		Text(w, http.StatusBadRequest, "notfqdn")
	default:
		Text(w, http.StatusInternalServerError, "911")
	}
}
