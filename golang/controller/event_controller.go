package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang/domain"
	"golang/service"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

const (
	eventNotFound       = "Event not found"
	notFound            = "Not found with ID"
	unexpected          = "Unexpected error"
	invalidUrl          = "Invalid ID in URL"
	internalServerError = "Internal server error"
)

type EventController struct {
	service service.IEventService
}

func NewEventController(service service.IEventService) *EventController {
	return &EventController{service: service}
}

func (c *EventController) Create(w http.ResponseWriter, r *http.Request) {
	var event domain.Event

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		slog.Info("Decode payload error", "error", err)
		http.Error(w, "Failed to Decode payload", http.StatusBadRequest)
		return
	}

	if err := event.Validate(); err != nil {
		slog.Info("Invalid Data", "error", err)
		http.Error(w, "Invalid Data", http.StatusBadRequest)
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
			slog.Error(unexpected, "error", err)
			http.Error(w, internalServerError, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", fmt.Sprintf("/events/%s", newEvent.ID))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newEvent)
}

func (c *EventController) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idFromPath := r.PathValue("id")

	if err := uuid.Validate(idFromPath); err != nil {
		slog.Info("Invalid ID", "error", err)
		http.Error(w, invalidUrl, http.StatusBadRequest)
		return
	}

	result, err := c.service.Get(ctx, idFromPath)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			slog.Info(notFound, "error", err)
			http.Error(w, eventNotFound, http.StatusNotFound)
		default:
			slog.Info(unexpected, "error", err)
			http.Error(w, internalServerError, http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(result)
}

func (c *EventController) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result, err := c.service.GetAll(ctx)

	if err != nil {
		slog.Info(unexpected, "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(result)
}

func (c *EventController) Update(w http.ResponseWriter, r *http.Request) {
	var updatedEvent domain.Event

	if err := json.NewDecoder(r.Body).Decode(&updatedEvent); err != nil {
		slog.Info("Decode payload error", "error", err)
		http.Error(w, "Failed to Decode payload", http.StatusBadRequest)
		return
	}

	idFromPath := r.PathValue("id")

	if err := uuid.Validate(idFromPath); err != nil {
		slog.Info("Invalid ID", "error", err)
		http.Error(w, invalidUrl, http.StatusBadRequest)
		return
	}

	updatedEvent.ID = idFromPath

	if err := c.service.Update(r.Context(), &updatedEvent); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			slog.Info(notFound, "error", err)
			http.Error(w, eventNotFound, http.StatusNotFound)
			return
		}
		http.Error(w, "Update failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(updatedEvent)
}

func (c *EventController) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idFromPath := r.PathValue("id")

	if err := uuid.Validate(idFromPath); err != nil {
		slog.Info("Invalid ID in URL", "error", err)
		http.Error(w, invalidUrl, http.StatusBadRequest)
		return
	}

	err := c.service.Delete(ctx, idFromPath)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			slog.Info(notFound, "error", err)
			http.Error(w, eventNotFound, http.StatusNotFound)
			return
		default:
			slog.Info(unexpected, "error", err)
			http.Error(w, internalServerError, http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
