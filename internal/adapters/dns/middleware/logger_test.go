package middleware

import (
	"github.com/miekg/dns"
	"go-dyndns/pkg/logger"
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

	mockLogger := logger.NewMockLogger(ctrl)

	msg := new(dns.Msg)
	msg.SetQuestion("example.com.", dns.TypeA)
	msg.Rcode = dns.RcodeSuccess

	mockLogger.EXPECT().Info(
		"DNS",
		"Request",
		logger.Field{Key: "domain", Value: "example.com."},
		logger.Field{Key: "type", Value: "A"},
		logger.Field{Key: "code", Value: "NOERROR"},
		logger.Field{Key: "client_ip", Value: "1.2.3.4"},
		gomock.Any(),
	)

	mockNextHandler := func(w dns.ResponseWriter, r *dns.Msg) {}
	handler := LoggingMiddleware(mockLogger, mockNextHandler)

	mockAddr := &net.TCPAddr{IP: net.ParseIP("1.2.3.4"), Port: 12345}
	writer := &mockResponseWriter{addr: mockAddr}

	handler(writer, msg)
}
