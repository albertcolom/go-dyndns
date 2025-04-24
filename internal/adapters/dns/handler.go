package dns

import (
	"context"
	"fmt"
	server "github.com/miekg/dns"
	"go-dyndns/internal/core/dns"
	"go-dyndns/pkg/logger"
	"strings"
)

type Handler struct {
	service dns.Service
	log     logger.Logger
}

func NewDnsHandler(service dns.Service, log logger.Logger) *Handler {
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
				h.log.Error(
					"DNS",
					"Error finding record",
					logger.Field{Key: "domain", Value: domainName},
					logger.Field{Key: "error", Value: err},
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
		h.log.Error("DNS", "Failed to write DNS response", logger.Field{Key: "error", Value: err})
	}
}
