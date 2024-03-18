package handlers

import (
	"context"
	"encoding/json"
	"filmography/internal/entities"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ActorService interface {
	CreateActor(ctx context.Context, actor entities.ActorEntity) error
	GetActors(ctx context.Context) ([]entities.ActorEntity, error)
	GetActor(ctx context.Context, id string) (entities.ActorEntity, error)
	UpdateActor(ctx context.Context, id string, actor entities.ActorEntity) error
	DeleteActor(ctx context.Context, id string) error
}

func (handlers Handlers) createActor(w http.ResponseWriter, r *http.Request) {
	actor := entities.ActorEntity{}
	err := json.NewDecoder(r.Body).Decode(&actor)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to decode JSON: %w", err).Error(), http.StatusBadRequest)
		logrus.WithField("error", err).Error("failed to decode JSON")
		return
	}

	err = handlers.svc.CreateActor(r.Context(), actor)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to create actor: %w", err).Error(), http.StatusInternalServerError)
		logrus.WithField("error", err).Error("failed to create actor")
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := map[string]string{
		"message": "actor is successfully created",
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

func (handlers Handlers) getActors(w http.ResponseWriter, r *http.Request) {
	actors, err := handlers.svc.GetActors(r.Context())
	if err != nil {
		http.Error(w, fmt.Errorf("failed to get actors: %w", err).Error(), http.StatusInternalServerError)
		logrus.WithField("error", err).Error("failed to get actors")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(actors)
	if err != nil {
		return
	}
}

func (handlers Handlers) getActor(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	actor, err := handlers.svc.GetActor(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to get actor: %w", err).Error(), http.StatusInternalServerError)
		logrus.WithField("error", err).Error("failed to get actor")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(actor)
	if err != nil {
		return
	}
}

func (handlers Handlers) updateActor(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	actor := entities.ActorEntity{}
	err := json.NewDecoder(r.Body).Decode(&actor)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to decode JSON: %w", err).Error(), http.StatusBadRequest)
		logrus.WithField("error", err).Error("failed to decode JSON")
		return
	}

	err = handlers.svc.UpdateActor(r.Context(), id, actor)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to update actor: %w", err).Error(), http.StatusInternalServerError)
		logrus.WithField("error", err).Error("failed to update actor")
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := map[string]string{
		"message": "actor is successfully updated",
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

func (handlers Handlers) deleteActor(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	err := handlers.svc.DeleteActor(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to delete actor: %w", err).Error(), http.StatusInternalServerError)
		logrus.WithField("error", err).Error("failed to delete actor")
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"message": "actor is successfully deleted",
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}
