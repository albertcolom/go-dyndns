package middleware

import (
	"context"
	"github.com/miekg/dns"
	"go-dyndns/internal/port"
	"net"
	"time"
)

func LoggingMiddleware(log port.Logger, next dns.HandlerFunc) dns.HandlerFunc {
	return func(w dns.ResponseWriter, r *dns.Msg) {
		start := time.Now()

		next(w, r)

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
			"code", dns.RcodeToString[r.Rcode],
			"client_ip", clientIP,
			"duration", time.Since(start),
		)
	}
}
