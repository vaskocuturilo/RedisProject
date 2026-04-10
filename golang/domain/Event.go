package domain

import (
	"errors"

	"github.com/google/uuid"
)

type Event struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func NewEvent(title, description string) *Event {
	return &Event{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
	}
}

var (
	ErrNotFound      = errors.New("event not found")
	ErrAlreadyExists = errors.New("event already exists")
	ErrInvalidInput  = errors.New("invalid input data")
)

func (e *Event) Validate() error {
	if e.Title == "" {
		return errors.New("title required")
	}

	if e.Description == "" {
		return errors.New("description required")
	}

	return nil
}
