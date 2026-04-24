package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang/domain"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

const eventString = "event:%s"

type CachedEventRepository struct {
	repo  EventRepository
	cache *redis.Client
	group singleflight.Group
	ttl   time.Duration
}

func NewCachedEventRepository(repo EventRepository, cache *redis.Client, ttl time.Duration) *CachedEventRepository {
	return &CachedEventRepository{
		repo:  repo,
		cache: cache,
		ttl:   ttl,
	}
}

func (c *CachedEventRepository) Create(ctx context.Context, event *domain.Event) error {
	if err := c.repo.Create(ctx, event); err != nil {
		return err
	}

	if event.ID == "" {
		return nil
	}

	data, err := json.Marshal(event)
	if err != nil {

		return err

	}
	key := fmt.Sprintf(eventString, event.ID)

	c.cache.Set(ctx, key, data, c.ttl)

	return nil
}

func (c *CachedEventRepository) Get(ctx context.Context, id string) (*domain.Event, error) {
	key := fmt.Sprintf(eventString, id)

	val, err := c.cache.Get(ctx, key).Result()

	if err == nil {
		var event domain.Event
		if err := json.Unmarshal([]byte(val), &event); err == nil {
			return &event, nil
		}
	}

	if err != nil && !errors.Is(err, redis.Nil) {
		slog.Info("Redis error for key ", "key", "error", key, err)
	}

	result, err, _ := c.group.Do(id, func() (interface{}, error) {
		dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		event, err := c.repo.Get(dbCtx, id)
		if err != nil {
			return nil, err
		}

		data, _ := json.Marshal(event)
		c.cache.Set(context.Background(), key, data, c.ttl)

		return event, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*domain.Event), nil
}

func (c *CachedEventRepository) GetAll(ctx context.Context) ([]*domain.Event, error) {
	return c.repo.GetAll(ctx)
}

func (c *CachedEventRepository) Update(ctx context.Context, event *domain.Event) error {
	key := fmt.Sprintf(eventString, event.ID)

	err := c.repo.Update(ctx, event)

	if err != nil {
		return err
	}

	if cacheErr := c.cache.Del(ctx, key).Err(); cacheErr != nil {
		slog.Info("failed to update cache key", "key", "error", key, cacheErr)
	}
	return nil
}

func (c *CachedEventRepository) Delete(ctx context.Context, id string) error {
	key := fmt.Sprintf(eventString, id)

	err := c.repo.Delete(ctx, id)

	if err != nil {
		return err
	}

	if cacheErr := c.cache.Del(ctx, key).Err(); cacheErr != nil {
		slog.Info("failed to delete cache key", "key", "error", key, cacheErr)
	}
	return nil
}
