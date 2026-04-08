package service

import (
	"context"
	"golang/domain"
	"golang/repository"
)

type EventService struct {
	repo repository.EventRepository
}

func NewEventService(repo repository.EventRepository) *EventService {
	return &EventService{repo: repo}
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
	//TODO implement me
	return nil, nil
}

func (s *EventService) Update(ctx context.Context, event *domain.Event) error {
	//TODO implement me
	return nil
}

func (s *EventService) Delete(ctx context.Context, id string) error {
	//TODO implement me
	return nil
}
