package handler

import "net/http"

func (h *Handler) Livez(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
