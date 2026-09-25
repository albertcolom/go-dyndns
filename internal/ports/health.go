//go:generate go tool mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

package ports

import "context"

// HealthChecker reports whether a dependency is reachable. Implementations
// must respect ctx and return promptly once it is done, since callers rely
// on it to bound the total time spent waiting on a readiness check.
type HealthChecker interface {
	Ping(ctx context.Context) error
}
