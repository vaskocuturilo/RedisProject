package repository

import (
	"context"
	"golang/domain"
)

type EventRepository interface {
	Create(ctx context.Context, event *domain.Event) error
	Get(ctx context.Context, id string) (*domain.Event, error)
	GetAll(ctx context.Context) ([]*domain.Event, error)
	Update(ctx context.Context, event *domain.Event) error
	Delete(ctx context.Context, id string) error
}
