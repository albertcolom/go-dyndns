package handler

import "net/http"

func (h *Handler) UpdateIp(w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Query().Get("domain")
	ip := r.URL.Query().Get("ip")

	if domain == "" || ip == "" {
		Error(w, http.StatusBadRequest, "Missing parameters")
		return
	}

	if err := h.service.Update(r.Context(), domain, ip); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	JSON(w, http.StatusOK, map[string]string{"message": "Updated " + domain + " to " + ip})
}

func (h *Handler) GetIp(w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Query().Get("domain")

	if domain == "" {
		Error(w, http.StatusBadRequest, "Missing parameters")
		return
	}

	record, err := h.service.Find(r.Context(), domain)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if record == nil {
		Error(w, http.StatusNotFound, "Domain not found")
		return
	}

	JSON(w, http.StatusOK, map[string]string{"domain": record.Domain, "ip": record.IP.String()})
}
