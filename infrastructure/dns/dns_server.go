package dns

import (
	"fmt"
	"log"

	"github.com/miekg/dns"
	"go-dyndns/domain"
)

type DNSServer struct {
	repository domain.DNSRepository
}

func NewDNSServer(repository domain.DNSRepository) *DNSServer {
	return &DNSServer{repository: repository}
}

func (s *DNSServer) Start() {
	dns.HandleFunc(".", s.handleDNSRequest)
}

func (s *DNSServer) handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	msg := new(dns.Msg)
	msg.SetReply(r)

	for _, question := range r.Question {
		if question.Qtype == dns.TypeA {
			domainName := question.Name

			record, err := s.repository.Find(domainName)
			if err != nil {
				log.Printf("Error finding record for %s: %v", domainName, err)
				continue
			}

			if record != nil {
				rr, err := dns.NewRR(fmt.Sprintf("%s A %s", domainName, record.IP.String()))
				if err == nil {
					msg.Answer = append(msg.Answer, rr)
				}
			}
		}
	}

	w.WriteMsg(msg)
}
