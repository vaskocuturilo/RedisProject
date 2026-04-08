package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"golang/domain"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockEventService struct {
	CreateFunc func(ctx context.Context, event *domain.Event) error
	GetFunc    func(ctx context.Context, id string) (*domain.Event, error)
	GetAllFunc func(ctx context.Context) ([]*domain.Event, error)
	UpdateFunc func(ctx context.Context, event *domain.Event) error
	DeleteFunc func(ctx context.Context, id string) error
}

func (m *MockEventService) Create(ctx context.Context, event *domain.Event) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, event)
	}
	return nil
}

func (m *MockEventService) Get(ctx context.Context, id string) (*domain.Event, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockEventService) GetAll(ctx context.Context) ([]*domain.Event, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx)
	}
	return nil, nil
}

func (m *MockEventService) Update(ctx context.Context, event *domain.Event) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, event)
	}
	return nil
}

func (m *MockEventService) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func TestEventController_Create_Success(t *testing.T) {
	// Arrange
	mockService := &MockEventService{
		CreateFunc: func(ctx context.Context, event *domain.Event) error {
			return nil
		},
	}
	ctrl := NewEventController(mockService)

	body, _ := json.Marshal(domain.Event{Title: "Test1", Description: "Test1"})

	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBuffer(body))

	rec := httptest.NewRecorder()

	// Act
	ctrl.Create(rec, req)

	// Assert
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, but got %d", rec.Code)
	}
}

func TestEventController_Create_ValidationAlreadyExist(t *testing.T) {
	// Arrange
	mockService := &MockEventService{
		CreateFunc: func(ctx context.Context, event *domain.Event) error {
			return domain.ErrAlreadyExists
		},
	}
	ctrl := NewEventController(mockService)

	body, _ := json.Marshal(domain.Event{Title: "Test1", Description: "Test1"})

	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBuffer(body))

	rec := httptest.NewRecorder()

	// Act
	ctrl.Create(rec, req)

	// Assert
	if rec.Code != http.StatusConflict {
		t.Errorf("Expected status 409 (Conflict), but got %d", rec.Code)
	}
}

func TestEventController_Create_ValidationError(t *testing.T) {
	// Arrange
	mockService := &MockEventService{
		CreateFunc: func(ctx context.Context, event *domain.Event) error {
			return domain.ErrInvalidInput
		},
	}
	ctrl := NewEventController(mockService)

	body, _ := json.Marshal(domain.Event{Title: "", Description: ""})

	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBuffer(body))

	rec := httptest.NewRecorder()

	// Act
	ctrl.Create(rec, req)

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rec.Code)
	}
}

func TestEventController_Create_ErrDecodeJSON(t *testing.T) {
	// Arrange
	mockService := &MockEventService{
		CreateFunc: func(ctx context.Context, event *domain.Event) error {
			return domain.ErrInvalidInput
		},
	}
	ctrl := NewEventController(mockService)

	body := []byte(`{"title": "Test", "description": "missing brace"`)

	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBuffer(body))

	rec := httptest.NewRecorder()

	// Act
	ctrl.Create(rec, req)

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for invalid JSON, got %d", rec.Code)
	}
}

func TestEventController_Create_ValidationInternalServerError(t *testing.T) {
	// Arrange
	mockService := &MockEventService{
		CreateFunc: func(ctx context.Context, event *domain.Event) error {
			return errors.New("something went wrong in the database")
		},
	}
	ctrl := NewEventController(mockService)

	body, _ := json.Marshal(domain.Event{Title: "test1", Description: "test2"})

	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBuffer(body))

	rec := httptest.NewRecorder()

	// Act
	ctrl.Create(rec, req)

	// Assert
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", rec.Code)
	}

	expectedBody := "Internal server error\n"
	if rec.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, rec.Body.String())
	}
}
