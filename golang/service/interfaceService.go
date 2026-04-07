package service

import (
	"context"
	"golang/domain"
)

type IEventService interface {
	Create(ctx context.Context, event *domain.Event) error
	GetEvent(ctx context.Context, id string) (*domain.Event, error)
	Get(ctx context.Context) ([]*domain.Event, error)
	Update(ctx context.Context, event *domain.Event) error
	Delete(ctx context.Context, id string) error
}
