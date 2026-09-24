package logger

import (
	"context"
	"log/slog"
	"os"

	"go-dyndns/internal/port"
)

type SlogLogger struct {
	slog *slog.Logger
}

var _ port.Logger = (*SlogLogger)(nil)

func NewSlogLogger(level string) *SlogLogger {
	var lvl slog.Level
	if err := lvl.UnmarshalText([]byte(level)); err != nil {
		lvl = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})
	return &SlogLogger{slog: slog.New(handler)}
}

func (l *SlogLogger) Debug(ctx context.Context, msg string, keysAndValues ...any) {
	l.slog.DebugContext(ctx, msg, keysAndValues...)
}

func (l *SlogLogger) Info(ctx context.Context, msg string, keysAndValues ...any) {
	l.slog.InfoContext(ctx, msg, keysAndValues...)
}

func (l *SlogLogger) Warn(ctx context.Context, msg string, keysAndValues ...any) {
	l.slog.WarnContext(ctx, msg, keysAndValues...)
}

func (l *SlogLogger) Error(ctx context.Context, msg string, keysAndValues ...any) {
	l.slog.ErrorContext(ctx, msg, keysAndValues...)
}

func (l *SlogLogger) With(keysAndValues ...any) port.Logger {
	return &SlogLogger{slog: l.slog.With(keysAndValues...)}
}
