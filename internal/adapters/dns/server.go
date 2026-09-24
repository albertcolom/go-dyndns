package dns

import (
	"context"
	"go-dyndns/internal/adapters/dns/middleware"
	"go-dyndns/internal/ports"

	server "github.com/miekg/dns"
)

type Server struct {
	DnsServer *server.Server
}

func NewDnsServer(handler *Handler, addr, net string, log ports.Logger) *Server {
	dnsServer := &server.Server{
		Addr: addr,
		Net:  net,
	}
	server.HandleFunc(".", middleware.LoggingMiddleware(log, handler.HandleDNSRequest))

	return &Server{DnsServer: dnsServer}
}

func (s *Server) Start() error {
	return s.DnsServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.DnsServer.ShutdownContext(ctx)
}
