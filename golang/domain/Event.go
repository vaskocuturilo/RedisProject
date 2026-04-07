package domain

import "errors"

type Event struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
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
