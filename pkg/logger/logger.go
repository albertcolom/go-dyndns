//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE

package logger

import "go.uber.org/zap"

type Logger interface {
	Info(context, msg string, fields ...Field)
	Error(context, msg string, fields ...Field)
	Debug(context, msg string, fields ...Field)
	Panic(context, msg string, fields ...Field)
	Fatal(context, msg string, fields ...Field)
	With(context string, fields ...Field) Logger
}

type Field struct {
	Key   string
	Value any
}

type ZapLogger struct {
	zap *zap.Logger
}

func NewZapLogger() (*ZapLogger, error) {
	z, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	z = z.WithOptions(zap.WithCaller(false), zap.AddCallerSkip(1))

	return &ZapLogger{zap: z}, nil
}

func (l *ZapLogger) Info(context, msg string, fields ...Field) {
	fields = append(fields, Field{Key: "context", Value: context})
	l.zap.Info(msg, toZapFields(fields...)...)
}

func (l *ZapLogger) Error(context, msg string, fields ...Field) {
	fields = append(fields, Field{Key: "context", Value: context})
	l.zap.Error(msg, toZapFields(fields...)...)
}

func (l *ZapLogger) Debug(context, msg string, fields ...Field) {
	fields = append(fields, Field{Key: "context", Value: context})
	l.zap.Debug(msg, toZapFields(fields...)...)
}

func (l *ZapLogger) Fatal(context, msg string, fields ...Field) {
	fields = append(fields, Field{Key: "context", Value: context})
	l.zap.Fatal(msg, toZapFields(fields...)...)
}

func (l *ZapLogger) Panic(context, msg string, fields ...Field) {
	fields = append(fields, Field{Key: "context", Value: context})
	l.zap.Fatal(msg, toZapFields(fields...)...)
}

func (l *ZapLogger) With(context string, fields ...Field) Logger {
	fields = append(fields, Field{Key: "context", Value: context})
	return &ZapLogger{zap: l.zap.With(toZapFields(fields...)...)}
}

func toZapFields(fields ...Field) []zap.Field {
	zf := make([]zap.Field, len(fields))
	for i, f := range fields {
		zf[i] = zap.Any(f.Key, f.Value)
	}
	return zf
}
