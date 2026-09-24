package http

import (
	"encoding/json"
	"net/http"

	"go-dyndns/internal/core/dns"
)

type Handler struct {
	service dns.Service
}

func NewHandler(service dns.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) UpdateIp(w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Query().Get("domain")
	ip := r.URL.Query().Get("ip")

	if domain == "" || ip == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing parameters"})
		return
	}

	if err := h.service.Update(r.Context(), domain, ip); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Updated " + domain + " to " + ip})
}

func (h *Handler) GetIp(w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Query().Get("domain")

	if domain == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing parameters"})
		return
	}

	record, err := h.service.Find(r.Context(), domain)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	if record == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Domain not found"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"domain": record.Domain, "ip": record.IP.String()})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
