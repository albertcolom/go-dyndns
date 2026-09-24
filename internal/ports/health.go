//go:generate go tool mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

package ports

import "context"

type HealthChecker interface {
	Ping(ctx context.Context) error
}
