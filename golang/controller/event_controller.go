package controller

import (
	"encoding/json"
	"errors"
	"fmt"
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

	newEvent := domain.NewEvent(event.Title, event.Description)

	err := c.service.Create(ctx, newEvent)

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

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", fmt.Sprintf("/events/%s", newEvent.ID))

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newEvent)
}
