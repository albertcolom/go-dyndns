package dns

import (
	"context"
	server "github.com/miekg/dns"
	"go-dyndns/internal/core/dns"
)

type Server struct {
	service dns.Service
	handler *Handler
}

func NewDnsServer(service dns.Service) *Server {
	handler := NewDnsHandler(service)
	return &Server{
		service: service,
		handler: handler,
	}
}

func (s *Server) Start(ctx context.Context, addr, net string) error {
	server.HandleFunc(".", func(w server.ResponseWriter, r *server.Msg) {
		s.handler.HandleDNSRequest(ctx, w, r)
	})

	dnsServer := &server.Server{Addr: addr, Net: net}

	return dnsServer.ListenAndServe()
}
