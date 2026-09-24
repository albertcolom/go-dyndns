package middleware

import (
	"context"
	"github.com/miekg/dns"
	"go-dyndns/internal/port"
	"net"
	"time"
)

type rcodeRecorder struct {
	dns.ResponseWriter
	rcode int
}

func (r *rcodeRecorder) WriteMsg(msg *dns.Msg) error {
	r.rcode = msg.Rcode
	return r.ResponseWriter.WriteMsg(msg)
}

func LoggingMiddleware(log port.Logger, next dns.HandlerFunc) dns.HandlerFunc {
	return func(w dns.ResponseWriter, r *dns.Msg) {
		start := time.Now()

		rec := &rcodeRecorder{ResponseWriter: w, rcode: dns.RcodeServerFailure}
		next(rec, r)

		var domain, qType string
		if len(r.Question) > 0 {
			domain = r.Question[0].Name
			qType = dns.TypeToString[r.Question[0].Qtype]
		}

		clientIP, _, _ := net.SplitHostPort(w.RemoteAddr().String())

		log.Info(context.Background(), "Request",
			"component", "DNS",
			"domain", domain,
			"type", qType,
			"code", dns.RcodeToString[rec.rcode],
			"client_ip", clientIP,
			"duration", time.Since(start),
		)
	}
}
