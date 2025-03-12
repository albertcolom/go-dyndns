package server

import (
	"fmt"
	"log"
	"strings"

	server "github.com/miekg/dns"
	"go-dyndns/internal/core/dns"
)

type Dns struct {
	service *dns.Service
}

func NewDns(service *dns.Service) *Dns {
	return &Dns{service: service}
}

func (s *Dns) Start(addr, net string) {
	server.HandleFunc(".", s.handleDNSRequest)

	server := &server.Server{Addr: addr, Net: net}

	log.Println("Starting DNS server on :53")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start DNS server: %v", err)
	}
}

func (s *Dns) handleDNSRequest(w server.ResponseWriter, r *server.Msg) {
	msg := new(server.Msg)
	msg.SetReply(r)

	for _, question := range r.Question {
		if question.Qtype == server.TypeA {
			domainName := strings.TrimSuffix(question.Name, ".")

			record, err := s.service.Find(domainName)
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

	w.WriteMsg(msg)
}
