package dns

import (
	"context"
	"fmt"
	"go-dyndns/internal/ports"
	"strings"

	server "github.com/miekg/dns"
)

type Handler struct {
	service ports.DNSService
	log     ports.Logger
}

func NewDnsHandler(service ports.DNSService, log ports.Logger) *Handler {
	return &Handler{service: service, log: log}
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
				h.log.Error(ctx, "Error finding record",
					"component", "DNS",
					"domain", domainName,
					"error", err,
				)
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
		h.log.Error(ctx, "Failed to write DNS response", "component", "DNS", "error", err)
	}
}
