package repository

import (
	"context"
	"database/sql"
	"errors"
	"golang/domain"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPostgresEventRepository_Create_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()

	defer db.Close()

	repo := NewPostgresEventRepository(db)

	event := &domain.Event{ID: "1", Title: "Title", Description: "Description"}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO events (id, title, description) VALUES ($1, $2, $3)")).
		WithArgs(event.ID, event.Title, event.Description).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Create(context.Background(), event)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestPostgresEventRepository_Update_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewPostgresEventRepository(db)

	event := &domain.Event{ID: "2", Title: "Title updated", Description: "Description updated"}

	mock.ExpectExec(regexp.QuoteMeta("UPDATE events SET title=$2, description=$3 WHERE id=$1")).
		WithArgs(event.ID, event.Title, event.Description).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Update(context.Background(), event)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostgresEventRepository_Delete_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewPostgresEventRepository(db)

	id := "1"

	mock.ExpectExec("DELETE").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Delete(context.Background(), id)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostgresEventRepository_Get_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPostgresEventRepository(db)

	id := "123"

	rows := sqlmock.NewRows([]string{"id", "title", "description"}).
		AddRow(id, "Test Title", "Test Description")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, description FROM events WHERE id = $1")).
		WithArgs(id).
		WillReturnRows(rows)

	event, err := repo.Get(context.Background(), id)

	if err != nil {
		t.Errorf("error was not expected: %s", err)
	}
	if event.Title != "Test Title" {
		t.Errorf("expected title 'Test Title', got %s", event.Title)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostgresEventRepository_Get_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewPostgresEventRepository(db)

	mock.ExpectQuery("SELECT").WillReturnError(sql.ErrNoRows)

	_, err := repo.Get(context.Background(), "unknown")

	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("expected domain.ErrNotFound, got %v", err)
	}
}

func TestPostgresEventRepository_GetAll_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPostgresEventRepository(db)

	expectedEvents := []domain.Event{
		{ID: "123", Title: "Test 1", Description: "Desc 1"},
		{ID: "456", Title: "Test 2", Description: "Desc 2"},
	}

	rows := sqlmock.NewRows([]string{"id", "title", "description"})

	for _, e := range expectedEvents {
		rows.AddRow(e.ID, e.Title, e.Description)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT *")).WillReturnRows(rows)

	events, err := repo.GetAll(context.Background())

	if len(events) != len(expectedEvents) {
		t.Fatalf("expected %d events, got %d", len(expectedEvents), len(events))
	}

	for i, e := range expectedEvents {
		if events[i].ID != e.ID {
			t.Errorf("[%d] expected ID %s, got %s", i, e.ID, events[i].ID)
		}

		if events[i].Title != e.Title {
			t.Errorf("[%d] expected Title %s, got %s", i, e.Title, events[i].Title)
		}

		if events[i].Description != e.Description {
			t.Errorf("[%d] expected Description %s, got %s", i, e.Description, events[i].Description)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %s", err)
	}
}

func TestPostgresEventRepository_GetAll_Empty(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewPostgresEventRepository(db)

	rows := sqlmock.NewRows([]string{"id", "title", "description"})

	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	events, err := repo.GetAll(context.Background())

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(events) != 0 {
		t.Errorf("expected 0 events, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %s", err)
	}
}
