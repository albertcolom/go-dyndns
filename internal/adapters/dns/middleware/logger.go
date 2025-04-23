package middleware

import (
	"github.com/miekg/dns"
	"log"
	"net"
	"time"
)

func LoggingMiddleware(next dns.HandlerFunc) dns.HandlerFunc {
	return func(w dns.ResponseWriter, r *dns.Msg) {
		start := time.Now()

		next(w, r)

		var domain, qType string
		if len(r.Question) > 0 {
			domain = r.Question[0].Name
			qType = dns.TypeToString[r.Question[0].Qtype]
		}

		duration := time.Since(start)
		clientIP, _, _ := net.SplitHostPort(w.RemoteAddr().String())
		rCode := dns.RcodeToString[r.Rcode]

		log.Printf("[DNS] %s IN %s | %s | %s | %s",
			domain, qType, rCode, duration, clientIP,
		)
	}
}
