package controller

import (
	"encoding/json"
	"errors"
	"golang/domain"
	"golang/service"
	"log"
	"net/http"
)

type EventController struct {
	service service.IEventService
}

func NewEventController(service service.IEventService) *EventController {
	return &EventController{service: service}
}

func (c *EventController) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	var event *domain.Event

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		log.Printf("Decode payload error: %v", err)
		http.Error(w, "Failed to Decode payload", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	err := c.service.Create(ctx, event)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidInput):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, domain.ErrAlreadyExists):
			http.Error(w, "Event ID already taken", http.StatusConflict)
		default:
			log.Printf("Unexpected error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
}
