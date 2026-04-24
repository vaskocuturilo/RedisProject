package service

import (
	"context"
	"fmt"
	"golang/domain"
	"golang/internal/config"
	"golang/repository"
	"log/slog"
)

type EventService struct {
	repo   repository.EventRepository
	locker repository.Locker
}

func NewEventService(repo repository.EventRepository, locker repository.Locker) *EventService {
	return &EventService{repo: repo, locker: locker}
}

func (s *EventService) Create(ctx context.Context, event *domain.Event) error {
	if err := event.Validate(); err != nil {
		return err
	}
	return s.repo.Create(ctx, event)
}

func (s *EventService) Get(ctx context.Context, id string) (*domain.Event, error) {
	return s.repo.Get(ctx, id)
}

func (s *EventService) GetAll(ctx context.Context) ([]*domain.Event, error) {
	return s.repo.GetAll(ctx)
}

func (s *EventService) Update(ctx context.Context, event *domain.Event) error {
	cfg := config.Load()

	lockKey := "event_update_" + event.ID

	lockValue, err := s.locker.Lock(ctx, lockKey, cfg.Server.RequestTimeout)

	if err != nil {
		return fmt.Errorf("resource is locked: %w", err)
	}

	defer func() {
		if unlockErr := s.locker.Unlock(ctx, lockKey, lockValue); unlockErr != nil {
			slog.Error("Failed to release lock", "key", lockKey, "error", unlockErr)
		}
	}()

	return s.repo.Update(ctx, event)
}

func (s *EventService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
