package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang/domain"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

const eventString = "event:%s"

type CachedEventRepository struct {
	repo  EventRepository
	cache *redis.Client
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

	data, _ := json.Marshal(event)
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
		log.Printf("Redis error for key %s: %v", key, err)
	}

	event, err := c.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(event)
	c.cache.Set(ctx, key, data, c.ttl)

	return event, nil
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
		log.Printf("failed to update cache key %s: %v", key, cacheErr)
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
		log.Printf("failed to delete cache key %s: %v", key, cacheErr)
	}
	return nil
}
