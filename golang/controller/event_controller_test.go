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

func TestEventController_Create_Table(t *testing.T) {
	tests := []struct {
		name       string
		body       any
		mockErr    error
		wantStatus int
	}{
		{
			name:       "Success",
			body:       domain.Event{Title: "Test1", Description: "Test1"},
			mockErr:    nil,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "Validation Already Exists",
			body:       domain.Event{Title: "Test1", Description: "Test1"},
			mockErr:    domain.ErrAlreadyExists,
			wantStatus: http.StatusConflict,
		},
		{
			name:       "Validation Error",
			body:       domain.Event{Title: "", Description: "Test1"},
			mockErr:    domain.ErrInvalidInput,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Validation Invalid JSON",
			body:       "not a json",
			mockErr:    nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Validation Internal Server Error",
			body:       domain.Event{Title: "Test1", Description: "Test1"},
			mockErr:    errors.New("something went wrong in the database"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService := &MockEventService{
				CreateFunc: func(ctx context.Context, e *domain.Event) error { return tc.mockErr },
			}
			ctrl := NewEventController(mockService)

			var buf bytes.Buffer
			json.NewEncoder(&buf).Encode(tc.body)

			req := httptest.NewRequest(http.MethodPost, "/events", &buf)
			rec := httptest.NewRecorder()

			ctrl.Create(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("got status %d, want %d", rec.Code, tc.wantStatus)
			}
		})
	}
}

func TestEventController_Get_Table(t *testing.T) {
	testEvent := &domain.Event{
		ID:          "550e8400-e29b-41d4-a716-446655440100",
		Title:       "Title",
		Description: "Description",
	}

	tests := []struct {
		name       string
		giveID     string
		wantReturn *domain.Event
		mockErr    error
		wantStatus int
	}{
		{
			name:       "Success",
			giveID:     testEvent.ID,
			wantReturn: testEvent,
			mockErr:    nil,
			wantStatus: http.StatusOK,
		},
		{
			name:       "Validation Incorrect ID",
			giveID:     "1",
			wantReturn: nil,
			mockErr:    domain.ErrInvalidInput,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Validation Not Found",
			giveID:     "550e8400-e29b-41d4-a716-446655440101",
			wantReturn: nil,
			mockErr:    domain.ErrNotFound,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Validation Internal Server Error",
			giveID:     testEvent.ID,
			wantReturn: nil,
			mockErr:    errors.New("something went wrong in the database"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService := &MockEventService{
				GetFunc: func(ctx context.Context, id string) (*domain.Event, error) {
					return tc.wantReturn, tc.mockErr
				},
			}
			ctrl := NewEventController(mockService)

			req := httptest.NewRequest(http.MethodGet, "/events/"+tc.giveID, nil)

			req.SetPathValue("id", tc.giveID)

			rec := httptest.NewRecorder()

			ctrl.Get(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("got status %d, want %d", rec.Code, tc.wantStatus)
			}

			if tc.wantStatus == http.StatusOK {
				var result domain.Event
				if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
					t.Fatalf("failed to unmarshal: %v", err)
				}
				if result.ID != tc.wantReturn.ID {
					t.Errorf("got ID %s, want %s", result.ID, tc.wantReturn.ID)
				}
			}
		})

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

func TestEventController_Update_Table(t *testing.T) {
	testEvent := &domain.Event{
		ID:          "550e8400-e29b-41d4-a716-446655440005",
		Title:       "Updated Title",
		Description: "Updated Description",
	}

	tests := []struct {
		name       string
		givenID    string
		body       any
		mockErr    error
		wantStatus int
	}{
		{
			name:       "Success",
			givenID:    testEvent.ID,
			body:       testEvent,
			mockErr:    nil,
			wantStatus: http.StatusOK,
		},
		{
			name:       "Validation Incorrect ID",
			givenID:    "1",
			body:       nil,
			mockErr:    domain.ErrInvalidInput,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Validation Not Found",
			givenID:    "550e8400-e29b-41d4-a716-446655440090",
			body:       nil,
			mockErr:    domain.ErrNotFound,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Validation Invalid JSON",
			givenID:    testEvent.ID,
			body:       "not a json",
			mockErr:    nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Validation Internal Server Error",
			givenID:    "550e8400-e29b-41d4-a716-446655440090",
			body:       nil,
			mockErr:    errors.New("something went wrong in the database"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService := &MockEventService{
				UpdateFunc: func(ctx context.Context, event *domain.Event) error {
					return tc.mockErr
				},
			}
			ctrl := NewEventController(mockService)

			body, _ := json.Marshal(tc.body)

			req := httptest.NewRequest(http.MethodPut, "/events/"+tc.givenID, bytes.NewBuffer(body))
			req.SetPathValue("id", tc.givenID)
			rec := httptest.NewRecorder()

			ctrl.Update(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("got status %d, want %d", rec.Code, tc.wantStatus)
			}
		})
	}
}

func TestEventController_Delete_Table(t *testing.T) {
	testEvent := &domain.Event{
		ID:          "550e8400-e29b-41d4-a716-446655440015",
		Title:       "Title",
		Description: "Description",
	}

	tests := []struct {
		name       string
		givenID    string
		body       any
		mockErr    error
		wantStatus int
	}{
		{
			name:       "Success",
			givenID:    testEvent.ID,
			body:       testEvent,
			mockErr:    nil,
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "Validation Incorrect ID",
			givenID:    "1",
			body:       nil,
			mockErr:    domain.ErrInvalidInput,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Validation Not Found",
			givenID:    "550e8400-e29b-41d4-a716-446655440090",
			body:       nil,
			mockErr:    domain.ErrNotFound,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Validation Internal Server Error",
			givenID:    "550e8400-e29b-41d4-a716-446655440090",
			body:       nil,
			mockErr:    errors.New("something went wrong in the database"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService := &MockEventService{
				DeleteFunc: func(ctx context.Context, id string) error {
					return tc.mockErr
				},
			}
			ctrl := NewEventController(mockService)

			body, _ := json.Marshal(tc.body)

			req := httptest.NewRequest(http.MethodDelete, "/events/"+tc.givenID, bytes.NewBuffer(body))
			req.SetPathValue("id", tc.givenID)
			rec := httptest.NewRecorder()

			ctrl.Delete(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("got status %d, want %d", rec.Code, tc.wantStatus)
			}
		})
	}
}
