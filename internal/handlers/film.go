package handlers

import (
	"context"
	"encoding/json"
	"filmography/internal/entities"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type FilmService interface {
	CreateFilm(ctx context.Context, film entities.FilmEntity) error
	GetFilms(ctx context.Context) ([]entities.FilmEntity, error)
	GetFilm(ctx context.Context, id string) (entities.FilmEntity, error)
	UpdateFilm(ctx context.Context, id string, film entities.FilmEntity) error
	DeleteFilm(ctx context.Context, id string) error
}

type CreateFilmRequest struct {
	Title       string                 `json:"name"`
	Description string                 `json:"description"`
	ReleaseDate time.Time              `json:"release_date"`
	Rating      float64                `json:"rating"`
	Actors      []entities.ActorEntity `json:"actors"`
}

// createFilm создает новый фильм.
// @Summary Создает фильм.
// @Description Создает новый фильм на основе переданных данных.
// @Tags Film
// @Accept json
// @Produce json
// @Param film body entities.FilmEntity true "Данные фильма"
// @Success 201 {object} map[string]string
// @Failure 400 {string} string "Ошибка при декодировании JSON"
// @Failure 500 {string} string "Ошибка при создании фильма"
// @Router /film [post]
func (handlers Handlers) createFilm(w http.ResponseWriter, r *http.Request) {
	request := CreateFilmRequest{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, fmt.Errorf("failed to decode JSON: %w", err).Error(), http.StatusBadRequest)
		logrus.WithField("error", err).Error("failed to decode JSON")
		return
	}

	if len(request.Title) == 0 || len(request.Title) > 150 {
		http.Error(w, "invalid name length", http.StatusBadRequest)
		logrus.Error("invalid name length")
		return
	}

	if len(request.Description) > 1000 {
		http.Error(w, "invalid description length", http.StatusBadRequest)
		logrus.Error("invalid description length")
		return
	}

	if request.Rating < 0 || request.Rating > 10 {
		http.Error(w, "invalid rating value", http.StatusBadRequest)
		logrus.Error("invalid rating value")
		return
	}

	film := entities.FilmEntity{
		Title:       request.Title,
		Description: request.Description,
		ReleaseDate: request.ReleaseDate,
		Rating:      request.Rating,
		Actors:      request.Actors,
	}

	err := handlers.svc.CreateFilm(r.Context(), film)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to create film: %w", err).Error(), http.StatusInternalServerError)
		logrus.WithField("error", err).Error("failed to create film")
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := map[string]string{
		"message": "film is successfully created",
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

// getFilms возвращает список фильмов.
// @Summary Возвращает список фильмов
// @Description Возвращает список всех фильмов.
// @Tags Film
// @Produce json
// @Success 200 {array} entities.FilmEntity "Список фильмов"
// @Failure 500 {string} string "Ошибка при получении фильмов"
// @Router /film [get]
func (handlers Handlers) getFilms(w http.ResponseWriter, r *http.Request) {
	films, err := handlers.svc.GetFilms(r.Context())
	if err != nil {
		http.Error(w, fmt.Errorf("failed to get films: %w", err).Error(), http.StatusInternalServerError)
		logrus.WithField("error", err).Error("failed to get films")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(films)
	if err != nil {
		return
	}
}

// getFilm возвращает информацию о фильме по его ID.
// @Summary Возвращает информацию о фильме
// @Description Возвращает информацию о фильме по указанному ID.
// @Tags Film
// @Param id query string true "ID фильма"
// @Produce json
// @Success 200 {object} entities.FilmEntity "Информация о фильме"
// @Failure 500 {string} string "Ошибка при получении фильма"
// @Router /film/{id} [get]
func (handlers Handlers) getFilm(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	film, err := handlers.svc.GetFilm(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to get film: %w", err).Error(), http.StatusInternalServerError)
		logrus.WithField("error", err).Error("failed to get film")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(film)
	if err != nil {
		return
	}
}

// updateFilm обновляет информацию о фильме.
// @Summary Обновляет информацию о фильме
// @Description Обновляет информацию о фильме с указанным ID на основе переданных данных.
// @Tags Film
// @Param id query string true "ID фильма"
// @Accept json
// @Produce json
// @Param film body entities.FilmEntity true "Данные фильма"
// @Success 201 {object} map[string]string
// @Failure 400 {string} string "Ошибка при декодировании JSON"
// @Failure 500 {string} string "Ошибка при обновлении фильма"
// @Router /film/{id} [put]
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

// deleteFilm удаляет фильма по его ID.
// @Summary Удаляет фильм
// @Description Удаляет фильм с указанным ID.
// @Tags Film
// @Param id query string true "ID фильма"
// @Success 200 {object} map[string]string
// @Failure 500 {string} string "Ошибка при удалении фильма"
// @Router /film/{id} [delete]
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
