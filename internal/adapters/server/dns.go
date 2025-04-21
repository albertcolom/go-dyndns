package server

import (
	"context"
	"fmt"
	"log"
	"strings"

	server "github.com/miekg/dns"
	"go-dyndns/internal/core/dns"
)

type Dns struct {
	service dns.Service
}

func NewDns(service dns.Service) *Dns {
	return &Dns{service: service}
}

func (s *Dns) Start(ctx context.Context, addr, net string) error {
	server.HandleFunc(".", func(w server.ResponseWriter, r *server.Msg) {
		s.handleDNSRequest(ctx, w, r)
	})

	dnsServer := &server.Server{Addr: addr, Net: net}

	return dnsServer.ListenAndServe()
}

func (s *Dns) handleDNSRequest(ctx context.Context, w server.ResponseWriter, r *server.Msg) {
	msg := new(server.Msg)
	msg.SetReply(r)

	for _, question := range r.Question {
		if question.Qtype == server.TypeA {
			domainName := strings.TrimSuffix(question.Name, ".")

			record, err := s.service.Find(ctx, domainName)
			if err != nil {
				log.Printf("Error finding record for %s: %v", domainName, err)
				continue
			}

			if record != nil {
				rr, err := server.NewRR(fmt.Sprintf("%s. A %s", domainName, record.IP.String()))
				if err == nil {
					msg.Answer = append(msg.Answer, rr)
				}
			}
		}
	}

	if err := w.WriteMsg(msg); err != nil {
		log.Printf("Failed to write DNS response: %v", err)
	}
}
