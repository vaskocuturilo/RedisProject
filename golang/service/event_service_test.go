package service

import (
	"context"
	"errors"
	"golang/domain"
	"golang/repository"
	"testing"
	"time"
)

type MockLocker struct {
	LockFunc   func(ctx context.Context, resource string, ttl time.Duration) (string, error)
	UnlockFunc func(ctx context.Context, resource string, lockValue string) error
}

func (m *MockLocker) Lock(ctx context.Context, res string, ttl time.Duration) (string, error) {
	if m.LockFunc != nil {
		return m.LockFunc(ctx, res, ttl)
	}
	return "default-lock-value", nil
}

func (m *MockLocker) Unlock(ctx context.Context, res string, val string) error {
	if m.UnlockFunc != nil {
		return m.UnlockFunc(ctx, res, val)
	}
	return nil
}

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
		giveTitle    string
		giveDesc     string
		mockResponse error
		mockLockErr  error
		wantErr      error
	}

	tests := []testCase{
		{
			name:         "Success",
			giveTitle:    "Valid Title",
			giveDesc:     "Valid Description",
			mockResponse: nil,
			mockLockErr:  nil,
			wantErr:      nil,
		},
		{
			name:         "Empty Title - Validation Error",
			giveTitle:    "",
			giveDesc:     "Some desc",
			mockResponse: nil,
			mockLockErr:  nil,
			wantErr:      domain.ErrTitleRequired,
		},
		{
			name:         "Repository Conflict",
			giveTitle:    "Duplicate",
			giveDesc:     "Desc",
			mockResponse: domain.ErrAlreadyExists,
			mockLockErr:  nil,
			wantErr:      domain.ErrAlreadyExists,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockRepository{
				CreateFunc: func(ctx context.Context, e *domain.Event) error {
					return tc.mockResponse
				},
			}

			mockLocker := &MockLocker{
				LockFunc: func(ctx context.Context, res string, ttl time.Duration) (string, error) {
					return "test-token", tc.mockLockErr
				},
			}
			serv := NewEventService(mockRepo, mockLocker)
			event := &domain.Event{Title: tc.giveTitle, Description: tc.giveDesc}

			// Act
			err := serv.Create(context.Background(), event)

			// Assert
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Expected error %v, but got %v", tc.wantErr, err)
			}
		})
	}
}

func TestEventService_Get_TableDriven(t *testing.T) {
	type testCase struct {
		name        string
		giveID      string
		wantReturn  *domain.Event
		mockErr     error
		mockLockErr error
		wantEvent   *domain.Event
		wantErr     error
	}

	event := domain.Event{Title: "Test", Description: "Test"}

	tests := []testCase{
		{
			name:        "Success",
			giveID:      "1",
			wantReturn:  &event,
			mockErr:     nil,
			mockLockErr: nil,
			wantEvent:   &event,
			wantErr:     nil,
		},
		{
			name:        "Not Found",
			giveID:      "2",
			wantReturn:  nil,
			mockErr:     domain.ErrNotFound,
			mockLockErr: nil,
			wantEvent:   nil,
			wantErr:     domain.ErrNotFound,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockRepository{
				GetFunc: func(ctx context.Context, id string) (*domain.Event, error) {
					if id != tc.giveID {
						t.Errorf("Mock expected ID %s, got %s", tc.giveID, id)
					}
					return tc.wantReturn, tc.mockErr
				},
			}

			mockLocker := &MockLocker{
				LockFunc: func(ctx context.Context, res string, ttl time.Duration) (string, error) {
					return "test-token", tc.mockLockErr
				},
			}

			serv := NewEventService(mockRepo, mockLocker)

			// Act
			result, err := serv.Get(context.Background(), tc.giveID)

			// Assert
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Expected error %v, got %v", tc.wantErr, err)
			}

			if tc.wantEvent != nil {
				if result == nil {
					t.Fatal("Expected result event, but got nil")
				}
				if result.ID != tc.wantEvent.ID {
					t.Errorf("Expected event ID %s, got %s", tc.wantEvent.ID, result.ID)
				}
			}
		})
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

	mockLocker := &MockLocker{
		LockFunc: func(ctx context.Context, res string, ttl time.Duration) (string, error) {
			return "test-token", nil
		},
	}

	//Act
	serv := NewEventService(mockRepo, mockLocker)

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

func TestEventService_Update_TableDriven(t *testing.T) {
	type testCase struct {
		name         string
		wantResponse error
		mockLockErr  error
		wantErr      error
	}

	event := domain.Event{ID: "1", Title: "Updated Test", Description: "Updated Description"}

	tests := []testCase{
		{
			name:         "Success",
			wantResponse: nil,
			mockLockErr:  nil,
			wantErr:      nil,
		},
		{
			name:         "Not Found",
			wantResponse: domain.ErrNotFound,
			mockLockErr:  nil,
			wantErr:      domain.ErrNotFound,
		},

		{
			name:         "Resource Locked",
			wantResponse: nil,
			mockLockErr:  repository.ErrLockNotAcquired,
			wantErr:      repository.ErrLockNotAcquired,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			repoCalled := false
			// Arrange
			mockRepo := &MockRepository{
				UpdateFunc: func(ctx context.Context, event *domain.Event) error {
					if event.ID != "1" {
						t.Errorf("Expected id '1', got '%s'", event.ID)
					}
					repoCalled = true
					return tc.wantResponse
				},
			}

			mockLocker := &MockLocker{
				LockFunc: func(ctx context.Context, res string, ttl time.Duration) (string, error) {
					return "test-token", tc.mockLockErr
				},
			}

			serv := NewEventService(mockRepo, mockLocker)

			// Act
			err := serv.Update(context.Background(), &event)

			// Assert
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Expected error %v, but got %v", tc.wantErr, err)
			}

			if tc.mockLockErr != nil && repoCalled {
				t.Error("Repository was called even though lock was not acquired!")
			}

			if tc.mockLockErr == nil && tc.name == "Success" && !repoCalled {
				t.Error("Repository was NOT called during successful lock acquisition")
			}
		})
	}
}

func TestEventService_Delete_TableDriven(t *testing.T) {
	type testCase struct {
		name         string
		giveID       string
		wantResponse error
		mockLockErr  error
		wantErr      error
	}

	tests := []testCase{
		{
			name:         "Success",
			giveID:       "1",
			wantResponse: nil,
			mockLockErr:  nil,
			wantErr:      nil,
		},
		{
			name:         "Not Found",
			giveID:       "1",
			wantResponse: domain.ErrNotFound,
			mockLockErr:  nil,
			wantErr:      domain.ErrNotFound,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockRepository{
				DeleteFunc: func(ctx context.Context, id string) error {
					if id != tc.giveID {
						t.Errorf("Expected id '1', got '%s'", id)
					}
					return tc.wantResponse
				},
			}

			mockLocker := &MockLocker{
				LockFunc: func(ctx context.Context, res string, ttl time.Duration) (string, error) {
					return "test-token", tc.mockLockErr
				},
			}

			serv := NewEventService(mockRepo, mockLocker)

			// Act
			err := serv.Delete(context.Background(), tc.giveID)

			// Assert
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Expected error %v, but got %v", tc.wantErr, err)
			}
		})
	}
}
