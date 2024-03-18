package handlers

import (
	"context"
	"encoding/json"
	"filmography/internal/entities"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

type FilmService interface {
	CreateFilm(ctx context.Context, film entities.FilmEntity) error
	GetFilms(ctx context.Context) ([]entities.FilmEntity, error)
	GetFilm(ctx context.Context, id string) (entities.FilmEntity, error)
	UpdateFilm(ctx context.Context, id string, film entities.FilmEntity) error
	DeleteFilm(ctx context.Context, id string) error
}

func (handlers Handlers) createFilm(w http.ResponseWriter, r *http.Request) {
	film := entities.FilmEntity{}
	err := json.NewDecoder(r.Body).Decode(&film)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to decode JSON: %w", err).Error(), http.StatusBadRequest)
		logrus.WithField("error", err).Error("failed to decode JSON")
		return
	}

	err = handlers.svc.CreateFilm(r.Context(), film)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to create film: %w", err).Error(), http.StatusInternalServerError)
		logrus.WithField("error", err).Error("failed to create film")
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

func (handlers Handlers) getFilms(w http.ResponseWriter, r *http.Request) {
	actors, err := handlers.svc.GetFilms(r.Context())
	if err != nil {
		http.Error(w, fmt.Errorf("failed to get films: %w", err).Error(), http.StatusInternalServerError)
		logrus.WithField("error", err).Error("failed to get films")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(actors)
	if err != nil {
		return
	}
}

func (handlers Handlers) getFilm(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	actor, err := handlers.svc.GetFilm(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to get film: %w", err).Error(), http.StatusInternalServerError)
		logrus.WithField("error", err).Error("failed to get film")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(actor)
	if err != nil {
		return
	}
}

func (handlers Handlers) updateFilm(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	film := entities.FilmEntity{}
	err := json.NewDecoder(r.Body).Decode(&film)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to decode JSON: %w", err).Error(), http.StatusBadRequest)
		logrus.WithField("error", err).Error("failed to decode JSON")
		return
	}

	err = handlers.svc.UpdateFilm(r.Context(), id, film)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to update film: %w", err).Error(), http.StatusInternalServerError)
		logrus.WithField("error", err).Error("failed to update film")
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := map[string]string{
		"message": "film is successfully updated",
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

func (handlers Handlers) deleteFilm(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	err := handlers.svc.DeleteFilm(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to delete film: %w", err).Error(), http.StatusInternalServerError)
		logrus.WithField("error", err).Error("failed to delete film")
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"message": "film is successfully deleted",
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}
