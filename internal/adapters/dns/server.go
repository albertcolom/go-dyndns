package dns

import (
	"context"
	server "github.com/miekg/dns"
	"go-dyndns/internal/adapters/dns/middleware"
	"go-dyndns/internal/core/dns"
)

type Server struct {
	DnsServer *server.Server
}

func NewDnsServer(service dns.Service, addr, net string) *Server {
	handler := NewDnsHandler(service)
	dnsServer := &server.Server{Addr: addr, Net: net}
	server.HandleFunc(".", middleware.LoggingMiddleware(handler.HandleDNSRequest))

	return &Server{DnsServer: dnsServer}

}

func (s *Server) Start() error {
	return s.DnsServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.DnsServer.ShutdownContext(ctx)
}
