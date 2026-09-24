//go:generate go tool mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

package port

import "context"

type Logger interface {
	Debug(ctx context.Context, msg string, keysAndValues ...any)
	Info(ctx context.Context, msg string, keysAndValues ...any)
	Warn(ctx context.Context, msg string, keysAndValues ...any)
	Error(ctx context.Context, msg string, keysAndValues ...any)
	With(keysAndValues ...any) Logger
}
