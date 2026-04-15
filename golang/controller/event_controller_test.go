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
	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status 201, but got %d", rec.Code)
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

func TestEventController_Get_Success(t *testing.T) {

	testEvent := &domain.Event{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		Title:       "Title",
		Description: "Description",
	}

	// Arrange
	mockService := &MockEventService{
		GetFunc: func(ctx context.Context, id string) (*domain.Event, error) {
			return testEvent, nil
		},
	}

	ctrl := NewEventController(mockService)

	req := httptest.NewRequest(http.MethodGet, "/events/"+testEvent.ID, nil)

	req.SetPathValue("id", testEvent.ID)

	rec := httptest.NewRecorder()

	// Act
	ctrl.Get(rec, req)

	// Assert
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, but got %d", rec.Code)
	}

	var result domain.Event

	err := json.Unmarshal(rec.Body.Bytes(), &result)

	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if result.ID != testEvent.ID {
		t.Errorf("Expected ID %s, but got %s", testEvent.ID, result.ID)
	}

	if result.Title != testEvent.Title {
		t.Errorf("Expected title %s, but got %s", testEvent.Title, result.Title)
	}
}

func TestEventController_Get_IncorrectID(t *testing.T) {
	expectedId := "1"

	// Arrange
	mockService := &MockEventService{
		GetFunc: func(ctx context.Context, id string) (*domain.Event, error) {
			return nil, domain.ErrInvalidInput
		},
	}

	ctrl := NewEventController(mockService)

	req := httptest.NewRequest(http.MethodGet, "/events/"+expectedId, nil)

	req.SetPathValue("id", expectedId)

	rec := httptest.NewRecorder()

	// Act
	ctrl.Get(rec, req)

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, but got %d", rec.Code)
	}
}

func TestEventController_Get_NotFound(t *testing.T) {
	expectedId := "550e8400-e29b-41d4-a716-446655440000"

	// Arrange
	mockService := &MockEventService{
		GetFunc: func(ctx context.Context, id string) (*domain.Event, error) {
			return nil, domain.ErrNotFound
		},
	}

	ctrl := NewEventController(mockService)

	req := httptest.NewRequest(http.MethodGet, "/events/"+expectedId, nil)

	req.SetPathValue("id", expectedId)

	rec := httptest.NewRecorder()

	// Act
	ctrl.Get(rec, req)

	// Assert
	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, but got %d", rec.Code)
	}
}

func TestEventController_Get_ValidationInternalServerError(t *testing.T) {
	expectedId := "550e8400-e29b-41d4-a716-446655440004"

	// Arrange
	mockService := &MockEventService{
		GetFunc: func(ctx context.Context, id string) (*domain.Event, error) {
			return nil, errors.New("something went wrong in the database")
		},
	}
	ctrl := NewEventController(mockService)

	req := httptest.NewRequest(http.MethodGet, "/events"+expectedId, nil)

	req.SetPathValue("id", expectedId)

	rec := httptest.NewRecorder()

	// Act
	ctrl.Get(rec, req)

	// Assert
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", rec.Code)
	}

	expectedBody := "Internal server error\n"
	if rec.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, rec.Body.String())
	}
}

func TestEventController_GetAll_Success(t *testing.T) {
	testEvents := []*domain.Event{
		{ID: "550e8400-e29b-41d4-a716-446655440001", Title: "Test 1", Description: "Desc 1"},
		{ID: "550e8400-e29b-41d4-a716-446655440002", Title: "Test 2", Description: "Desc 2"},
	}

	// Arrange
	mockService := &MockEventService{
		GetAllFunc: func(ctx context.Context) ([]*domain.Event, error) {
			return testEvents, nil
		},
	}

	ctrl := NewEventController(mockService)

	req := httptest.NewRequest(http.MethodGet, "/events/", nil)

	rec := httptest.NewRecorder()

	// Act
	ctrl.GetAll(rec, req)

	// Assert
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, but got %d", rec.Code)
	}

	var result []*domain.Event

	err := json.Unmarshal(rec.Body.Bytes(), &result)

	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(result) != len(testEvents) {
		t.Errorf("Expected len of testEvents %d, is equal len(result) %d", len(testEvents), len(result))
	}
}

func TestEventController_Update_Success(t *testing.T) {
	testEvent := &domain.Event{
		ID:          "550e8400-e29b-41d4-a716-446655440005",
		Title:       "Updated Title",
		Description: "Updated Description",
	}

	// Arrange
	mockService := &MockEventService{
		UpdateFunc: func(ctx context.Context, event *domain.Event) error {
			return nil
		},
	}

	body, _ := json.Marshal(testEvent)

	ctrl := NewEventController(mockService)

	req := httptest.NewRequest(http.MethodPut, "/events/"+testEvent.ID, bytes.NewBuffer(body))

	req.SetPathValue("id", testEvent.ID)

	rec := httptest.NewRecorder()

	// Act
	ctrl.Update(rec, req)

	// Assert
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, but got %d", rec.Code)
	}

	var result domain.Event

	err := json.Unmarshal(rec.Body.Bytes(), &result)

	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if result.ID != testEvent.ID {
		t.Errorf("Expected ID %s, but got %s", testEvent.ID, result.ID)
	}

	if result.Title != testEvent.Title {
		t.Errorf("Expected title %s, but got %s", testEvent.Title, result.Title)
	}
}

func TestEventController_Update_IncorrectID(t *testing.T) {
	testEvent := &domain.Event{
		ID:          "1",
		Title:       "Updated Title",
		Description: "Updated Description",
	}

	// Arrange
	mockService := &MockEventService{
		UpdateFunc: func(ctx context.Context, event *domain.Event) error {
			return domain.ErrInvalidInput
		},
	}

	body, _ := json.Marshal(testEvent)

	ctrl := NewEventController(mockService)

	req := httptest.NewRequest(http.MethodPut, "/events/"+testEvent.ID, bytes.NewBuffer(body))

	req.SetPathValue("id", testEvent.ID)

	rec := httptest.NewRecorder()

	// Act
	ctrl.Update(rec, req)

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, but got %d", rec.Code)
	}
}

func TestEventController_Update_NotFound(t *testing.T) {
	testEvent := &domain.Event{
		ID:          "550e8400-e29b-41d4-a716-446655440010",
		Title:       "Updated Title",
		Description: "Updated Description",
	}

	// Arrange
	mockService := &MockEventService{
		UpdateFunc: func(ctx context.Context, event *domain.Event) error {
			return domain.ErrNotFound
		},
	}

	body, _ := json.Marshal(testEvent)

	ctrl := NewEventController(mockService)

	req := httptest.NewRequest(http.MethodPut, "/events/"+testEvent.ID, bytes.NewBuffer(body))

	req.SetPathValue("id", testEvent.ID)

	rec := httptest.NewRecorder()

	// Act
	ctrl.Update(rec, req)

	// Assert
	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, but got %d", rec.Code)
	}
}

func TestEventController_Update_ErrDecodeJSON(t *testing.T) {
	expectedId := "550e8400-e29b-41d4-a716-446655440014"

	// Arrange
	mockService := &MockEventService{
		UpdateFunc: func(ctx context.Context, event *domain.Event) error {
			return domain.ErrInvalidInput
		},
	}
	ctrl := NewEventController(mockService)

	body := []byte(`{"title": "Test", "description": "missing brace"`)

	req := httptest.NewRequest(http.MethodPut, "/events"+expectedId, bytes.NewBuffer(body))

	req.SetPathValue("id", expectedId)

	rec := httptest.NewRecorder()

	// Act
	ctrl.Update(rec, req)

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for invalid JSON, got %d", rec.Code)
	}
}

func TestEventController_Update_ValidationInternalServerError(t *testing.T) {
	testEvent := &domain.Event{
		ID:          "550e8400-e29b-41d4-a716-446655440010",
		Title:       "Updated Title",
		Description: "Updated Description",
	}

	// Arrange
	mockService := &MockEventService{
		UpdateFunc: func(ctx context.Context, event *domain.Event) error {
			return errors.New("something went wrong in the database")
		},
	}

	body, _ := json.Marshal(testEvent)

	ctrl := NewEventController(mockService)

	req := httptest.NewRequest(http.MethodPut, "/events/"+testEvent.ID, bytes.NewBuffer(body))

	req.SetPathValue("id", testEvent.ID)

	rec := httptest.NewRecorder()

	// Act
	ctrl.Update(rec, req)

	// Assert
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, but got %d", rec.Code)
	}
}

func TestEventController_Delete_Success(t *testing.T) {
	expectedId := "550e8400-e29b-41d4-a716-446655440300"

	// Arrange
	mockService := &MockEventService{
		DeleteFunc: func(ctx context.Context, id string) error {
			return nil
		},
	}

	ctrl := NewEventController(mockService)

	req := httptest.NewRequest(http.MethodDelete, "/events/"+expectedId, nil)

	req.SetPathValue("id", expectedId)

	rec := httptest.NewRecorder()

	// Act
	ctrl.Delete(rec, req)

	// Assert
	if rec.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, but got %d", rec.Code)
	}
}

func TestEventController_Delete_IncorrectID(t *testing.T) {
	expectedId := "1"

	// Arrange
	mockService := &MockEventService{
		DeleteFunc: func(ctx context.Context, id string) error {
			return domain.ErrInvalidInput
		},
	}

	ctrl := NewEventController(mockService)

	req := httptest.NewRequest(http.MethodDelete, "/events/"+expectedId, nil)

	req.SetPathValue("id", expectedId)

	rec := httptest.NewRecorder()

	// Act
	ctrl.Delete(rec, req)

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, but got %d", rec.Code)
	}
}

func TestEventController_Delete_NotFound(t *testing.T) {
	expectedId := "550e8400-e29b-41d4-a716-446655440300"

	// Arrange
	mockService := &MockEventService{
		DeleteFunc: func(ctx context.Context, id string) error {
			return domain.ErrNotFound
		},
	}

	ctrl := NewEventController(mockService)

	req := httptest.NewRequest(http.MethodDelete, "/events/"+expectedId, nil)

	req.SetPathValue("id", expectedId)

	rec := httptest.NewRecorder()

	// Act
	ctrl.Delete(rec, req)

	// Assert
	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, but got %d", rec.Code)
	}
}

func TestEventController_Delete_ValidationInternalServerError(t *testing.T) {
	expectedId := "550e8400-e29b-41d4-a716-446655440300"

	// Arrange
	mockService := &MockEventService{
		DeleteFunc: func(ctx context.Context, id string) error {
			return errors.New("something went wrong in the database")
		},
	}

	ctrl := NewEventController(mockService)

	req := httptest.NewRequest(http.MethodDelete, "/events/"+expectedId, nil)

	req.SetPathValue("id", expectedId)

	rec := httptest.NewRecorder()

	// Act
	ctrl.Delete(rec, req)

	// Assert
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, but got %d", rec.Code)
	}
}
