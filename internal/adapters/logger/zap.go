package logger

import (
	"context"

	"go-dyndns/internal/port"

	"go.uber.org/zap"
)

type ZapLogger struct {
	zap *zap.SugaredLogger
}

var _ port.Logger = (*ZapLogger)(nil)

func NewZapLogger() (*ZapLogger, error) {
	z, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	z = z.WithOptions(zap.WithCaller(false), zap.AddCallerSkip(1))

	return &ZapLogger{zap: z.Sugar()}, nil
}

func (l *ZapLogger) Debug(_ context.Context, msg string, keysAndValues ...any) {
	l.zap.Debugw(msg, keysAndValues...)
}

func (l *ZapLogger) Info(_ context.Context, msg string, keysAndValues ...any) {
	l.zap.Infow(msg, keysAndValues...)
}

func (l *ZapLogger) Warn(_ context.Context, msg string, keysAndValues ...any) {
	l.zap.Warnw(msg, keysAndValues...)
}

func (l *ZapLogger) Error(_ context.Context, msg string, keysAndValues ...any) {
	l.zap.Errorw(msg, keysAndValues...)
}

func (l *ZapLogger) With(keysAndValues ...any) port.Logger {
	return &ZapLogger{zap: l.zap.With(keysAndValues...)}
}
