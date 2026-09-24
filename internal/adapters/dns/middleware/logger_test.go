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

func (m *mockResponseWriter) WriteMsg(msg *dns.Msg) error {
	return nil
}

func TestLoggingMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := port.NewMockLogger(ctrl)

	msg := new(dns.Msg)
	msg.SetQuestion("example.com.", dns.TypeA)

	mockLogger.EXPECT().Info(
		gomock.Any(),
		"Request",
		"component", "DNS",
		"domain", "example.com.",
		"type", "A",
		"code", "NXDOMAIN",
		"client_ip", "1.2.3.4",
		"duration", gomock.Any(),
	)

	// next writes a reply whose Rcode differs from the request's, proving the
	// middleware logs the actual response code and not the request's.
	mockNextHandler := func(w dns.ResponseWriter, r *dns.Msg) {
		reply := new(dns.Msg)
		reply.SetReply(r)
		reply.Rcode = dns.RcodeNameError
		_ = w.WriteMsg(reply)
	}
	handler := LoggingMiddleware(mockLogger, mockNextHandler)

	mockAddr := &net.TCPAddr{IP: net.ParseIP("1.2.3.4"), Port: 12345}
	writer := &mockResponseWriter{addr: mockAddr}

	handler(writer, msg)
}
