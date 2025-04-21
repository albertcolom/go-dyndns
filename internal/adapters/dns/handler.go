package dns

import (
	"context"
	"fmt"
	server "github.com/miekg/dns"
	"go-dyndns/internal/core/dns"
	"log"
	"strings"
)

type Handler struct {
	service dns.Service
}

func NewDnsHandler(service dns.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) HandleDNSRequest(w server.ResponseWriter, r *server.Msg) {
	ctx := context.Background()
	msg := new(server.Msg)
	msg.SetReply(r)

	for _, question := range r.Question {
		if question.Qtype == server.TypeA {
			domainName := strings.TrimSuffix(question.Name, ".")
			record, err := h.service.Find(ctx, domainName)
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
