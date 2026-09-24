package middleware

import (
	"github.com/miekg/dns"
	"go-dyndns/internal/port"
	"go.uber.org/mock/gomock"
	"net"
	"testing"
)

type mockResponseWriter struct {
	addr net.Addr
	dns.ResponseWriter
}

func (m *mockResponseWriter) RemoteAddr() net.Addr {
	return m.addr
}

func TestLoggingMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := port.NewMockLogger(ctrl)

	msg := new(dns.Msg)
	msg.SetQuestion("example.com.", dns.TypeA)
	msg.Rcode = dns.RcodeSuccess

	mockLogger.EXPECT().Info(
		gomock.Any(),
		"Request",
		"component", "DNS",
		"domain", "example.com.",
		"type", "A",
		"code", "NOERROR",
		"client_ip", "1.2.3.4",
		"duration", gomock.Any(),
	)

	mockNextHandler := func(w dns.ResponseWriter, r *dns.Msg) {}
	handler := LoggingMiddleware(mockLogger, mockNextHandler)

	mockAddr := &net.TCPAddr{IP: net.ParseIP("1.2.3.4"), Port: 12345}
	writer := &mockResponseWriter{addr: mockAddr}

	handler(writer, msg)
}
