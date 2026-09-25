package dns

import (
	"context"
	"fmt"
	"net"
	"time"

	"go-dyndns/internal/adapters/dns/middleware"
	"go-dyndns/internal/ports"

	server "github.com/miekg/dns"
)

const pingTimeout = 2 * time.Second

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

func (s *Server) Ping(ctx context.Context) error {
	host, port, err := net.SplitHostPort(s.DnsServer.Addr)
	if err != nil {
		return fmt.Errorf("invalid DNS server address %q: %w", s.DnsServer.Addr, err)
	}

	msg := new(server.Msg)
	msg.SetQuestion("healthcheck.go-dyndns.internal.", server.TypeA)

	client := &server.Client{Net: s.DnsServer.Net, Timeout: pingTimeout}
	if _, _, err := client.ExchangeContext(ctx, msg, net.JoinHostPort(host, port)); err != nil {
		return fmt.Errorf("DNS self-check failed: %w", err)
	}

	return nil
}
