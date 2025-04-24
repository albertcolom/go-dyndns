package middleware

import (
	"github.com/miekg/dns"
	"go-dyndns/pkg/logger"
	"net"
	"time"
)

func LoggingMiddleware(log logger.Logger, next dns.HandlerFunc) dns.HandlerFunc {
	return func(w dns.ResponseWriter, r *dns.Msg) {
		start := time.Now()

		next(w, r)

		var domain, qType string
		if len(r.Question) > 0 {
			domain = r.Question[0].Name
			qType = dns.TypeToString[r.Question[0].Qtype]
		}

		clientIP, _, _ := net.SplitHostPort(w.RemoteAddr().String())

		log.Info(
			"DNS",
			"Request",
			logger.Field{Key: "domain", Value: domain},
			logger.Field{Key: "type", Value: qType},
			logger.Field{Key: "code", Value: dns.RcodeToString[r.Rcode]},
			logger.Field{Key: "client_ip", Value: clientIP},
			logger.Field{Key: "duration", Value: time.Since(start)},
		)
	}
}
