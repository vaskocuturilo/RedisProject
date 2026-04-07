package repository

import (
	"context"
	"golang/domain"
)

type EventRepository interface {
	Create(ctx context.Context, event *domain.Event) error
	FindByID(ctx context.Context, id string) (*domain.Event, error)
	FindAll(ctx context.Context) ([]*domain.Event, error)
	Update(ctx context.Context, event *domain.Event) error
	Delete(ctx context.Context, id string) error
}
