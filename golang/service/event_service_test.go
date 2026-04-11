package service

import (
	"context"
	"errors"
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
	if m.GetFunc != nil {

		return m.GetFunc(ctx, id)

	}
	return nil, nil
}
func (m *MockRepository) GetAll(ctx context.Context) ([]*domain.Event, error) {
	if m.GetAllFunc != nil {

		return m.GetAllFunc(ctx)

	}
	return nil, nil
}
func (m *MockRepository) Update(ctx context.Context, event *domain.Event) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, event)
	}
	return nil
}
func (m *MockRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func TestEventService_Create_TableDriven(t *testing.T) {
	type testCase struct {
		name         string
		inputTitle   string
		inputDesc    string
		mockResponse error
		wantErr      error
	}

	tests := []testCase{
		{
			name:         "Success",
			inputTitle:   "Valid Title",
			inputDesc:    "Valid Description",
			mockResponse: nil,
			wantErr:      nil,
		},
		{
			name:         "Empty Title - Validation Error",
			inputTitle:   "",
			inputDesc:    "Some desc",
			mockResponse: nil,
			wantErr:      domain.ErrTitleRequired,
		},
		{
			name:         "Repository Conflict",
			inputTitle:   "Duplicate",
			inputDesc:    "Desc",
			mockResponse: domain.ErrAlreadyExists,
			wantErr:      domain.ErrAlreadyExists,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockRepository{
				CreateFunc: func(ctx context.Context, e *domain.Event) error {
					return tc.mockResponse
				},
			}
			serv := NewEventService(mockRepo)
			event := &domain.Event{Title: tc.inputTitle, Description: tc.inputDesc}

			// Act
			err := serv.Create(context.Background(), event)

			// Assert
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Expected error %v, but got %v", tc.wantErr, err)
			}
		})
	}
}

func TestEventService_Get_Success(t *testing.T) {
	//Arrange
	event := domain.Event{ID: "1", Title: "Test", Description: "Test"}

	mockRepo := &MockRepository{
		GetFunc: func(ctx context.Context, id string) (*domain.Event, error) {
			return &event, nil
		},
	}
	//Act
	serv := NewEventService(mockRepo)

	ctx := context.Background()

	result, err := serv.Get(ctx, event.ID)

	//Assert
	if err != nil {
		t.Fatalf("Expected success, got err: %v", err)
	}

	if result.ID != event.ID {
		t.Errorf("Expected ID result ID: %v is equal event ID: %v, got err: %v", result.ID, event.ID, err)
	}
}

func TestEventService_Get_NotFound(t *testing.T) {
	//Arrange
	mockRepo := &MockRepository{
		GetFunc: func(ctx context.Context, id string) (*domain.Event, error) {
			return nil, domain.ErrNotFound
		},
	}
	//Act
	serv := NewEventService(mockRepo)

	ctx := context.Background()

	_, err := serv.Get(ctx, "2")

	//Assert
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("Expected event not found error, but got nil")
	}
}

func TestEventService_GetAll_Success(t *testing.T) {
	//Arrange
	expectedEvents := []*domain.Event{
		{ID: "123", Title: "Test 1", Description: "Desc 1"},
		{ID: "456", Title: "Test 2", Description: "Desc 2"},
	}

	mockRepo := &MockRepository{
		GetAllFunc: func(ctx context.Context) ([]*domain.Event, error) {
			return expectedEvents, nil
		},
	}
	//Act
	serv := NewEventService(mockRepo)

	ctx := context.Background()

	result, err := serv.GetAll(ctx)

	//Assert
	if err != nil {
		t.Fatalf("Expected success, got err: %v", err)
	}

	if len(result) != len(expectedEvents) {
		t.Errorf("Expected len result: %v is equal len expectedEvents: %v, got err: %v", len(result), len(expectedEvents), err)
	}
}

func TestEventService_Update_Success(t *testing.T) {
	//Arrange
	event := domain.Event{ID: "1", Title: "Updated Test", Description: "Updated Description"}

	mockRepo := &MockRepository{
		UpdateFunc: func(ctx context.Context, event *domain.Event) error {
			if event.ID != "1" {
				t.Errorf("Expected id '1', got '%s'", event.ID)
			}
			return nil
		},
	}
	//Act
	serv := NewEventService(mockRepo)

	ctx := context.Background()

	err := serv.Update(ctx, &event)

	//Assert
	if err != nil {
		t.Fatalf("Expected success, got err: %v", err)
	}
}

func TestEventService_Update_NotFound(t *testing.T) {
	//Arrange
	event := domain.Event{ID: "", Title: "", Description: ""}

	mockRepo := &MockRepository{
		UpdateFunc: func(ctx context.Context, event *domain.Event) error {
			return domain.ErrNotFound
		},
	}
	//Act
	serv := NewEventService(mockRepo)

	ctx := context.Background()

	err := serv.Update(ctx, &event)

	//Assert
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("Expected event not found error, but got nil")
	}
}

func TestEventService_Delete_Success(t *testing.T) {
	//Arrange
	mockRepo := &MockRepository{
		DeleteFunc: func(ctx context.Context, id string) error {
			if id != "1" {
				t.Errorf("Expected id '1', got '%s'", id)
			}
			return nil
		},
	}
	//Act
	serv := NewEventService(mockRepo)

	ctx := context.Background()

	err := serv.Delete(ctx, "1")

	//Assert
	if err != nil {
		t.Fatalf("Expected success, got err: %v", err)
	}
}

func TestEventService_Delete_NotFound(t *testing.T) {
	//Arrange
	mockRepo := &MockRepository{
		DeleteFunc: func(ctx context.Context, id string) error {
			if id != "1" {
				t.Errorf("Expected id '1', got '%s'", id)
			}
			return domain.ErrNotFound
		},
	}
	//Act
	serv := NewEventService(mockRepo)

	ctx := context.Background()

	err := serv.Delete(ctx, "1")

	//Assert
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("Expected event not found error, but got nil")
	}
}
