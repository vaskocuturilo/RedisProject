package service

import (
	"context"
	"golang/domain"
	"testing"
)

type MockRepository struct {
	CreateFunc func(ctx context.Context, event *domain.Event) error
	GetFunc    func(ctx context.Context, id string) (*domain.Event, error)
	GetAllFunc func(ctx context.Context) ([]*domain.Event, error)
	UpdateFunc func(ctx context.Context, event *domain.Event) error
	DeleteFunc func(ctx context.Context, id string) error
}

func (m *MockRepository) Create(ctx context.Context, event *domain.Event) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, event)
	}
	return nil
}

func (m *MockRepository) Get(ctx context.Context, id string) (*domain.Event, error) {
	return nil, nil
}
func (m *MockRepository) GetAll(ctx context.Context) ([]*domain.Event, error)   { return nil, nil }
func (m *MockRepository) Update(ctx context.Context, event *domain.Event) error { return nil }
func (m *MockRepository) Delete(ctx context.Context, id string) error           { return nil }

func TestEventService_Success(t *testing.T) {
	//Arrange
	event := domain.Event{Title: "Test", Description: "Test"}

	wasCalled := false

	mockRepo := &MockRepository{
		CreateFunc: func(ctx context.Context, event *domain.Event) error {
			wasCalled = true
			return nil
		},
	}
	//Act
	serv := NewEventService(mockRepo)

	ctx := context.Background()

	err := serv.Create(ctx, &event)

	//Assert
	if err != nil {
		t.Errorf("Expected success, got err: %v", err)
	}
	if !wasCalled {
		t.Error("Repository was not called, but it should have been!")
	}
}

func TestEventService_ValidationError(t *testing.T) {
	//Arrange
	event := domain.Event{Title: "", Description: "Test"}

	wasCalled := false

	mockRepo := &MockRepository{
		CreateFunc: func(ctx context.Context, event *domain.Event) error {
			wasCalled = true
			return nil
		},
	}
	//Act
	serv := NewEventService(mockRepo)

	ctx := context.Background()

	err := serv.Create(ctx, &event)

	//Assert
	if err == nil {
		t.Errorf("Expected validation error, but got nil")
	}
	if wasCalled {
		t.Error("Repository was not called, but it should have been!")
	}
}
