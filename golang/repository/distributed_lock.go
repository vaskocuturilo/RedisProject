package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var ErrLockNotAcquired = errors.New("could not acquire lock")

type RedisLock struct {
	client *redis.Client
}

func NewRedisLock(client *redis.Client) *RedisLock {
	return &RedisLock{client: client}
}

var unlockScript = redis.NewScript(`
    if redis.call("GET", KEYS[1]) == ARGV[1] then

    return redis.call("DEL", KEYS[1])

else

    return 0

end
`)

func (l *RedisLock) Lock(ctx context.Context, resource string, ttl time.Duration) (string, error) {
	lockValue := uuid.New().String()

	success, err := l.client.SetNX(ctx, "lock:"+resource, lockValue, ttl).Result()
	if err != nil {
		return "", err
	}
	if !success {
		return "", ErrLockNotAcquired
	}

	return lockValue, nil
}

func (l *RedisLock) Unlock(ctx context.Context, resource string, lockValue string) error {
	return unlockScript.Run(ctx, l.client, []string{"lock:" + resource}, lockValue).Err()
}
