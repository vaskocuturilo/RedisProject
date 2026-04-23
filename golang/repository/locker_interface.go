package repository

import (
	"context"
	"time"
)

type Locker interface {
	Lock(ctx context.Context, resource string, ttl time.Duration) (string, error)
	Unlock(ctx context.Context, resource string, lockValue string) error
}
