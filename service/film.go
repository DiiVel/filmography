package service

import (
	"context"
	"filmography/internal/entities"
	"fmt"
	"github.com/google/uuid"
)

type FilmService struct {
	repo FilmRepoInterface
}

type FilmRepoInterface interface {
	CreateFilm(ctx context.Context, film entities.FilmEntity) error
	GetFilms(ctx context.Context) ([]entities.FilmEntity, error)
	GetFilm(ctx context.Context, id string) (entities.FilmEntity, error)
	UpdateFilm(ctx context.Context, id string, film entities.FilmEntity) error
	DeleteFilm(ctx context.Context, id string) error
}

func NewFilmService(repo FilmRepoInterface) FilmService {
	return FilmService{
		repo: repo,
	}
}

func (svc FilmService) CreateFilm(ctx context.Context, film entities.FilmEntity) error {
	film.ID = uuid.NewString()
	err := svc.repo.CreateFilm(ctx, film)
	return err
}

func (svc FilmService) GetFilms(ctx context.Context) ([]entities.FilmEntity, error) {
	films, err := svc.repo.GetFilms(ctx)
	if err != nil {
		return nil, fmt.Errorf("get films failed: %w", err)
	}
	return films, err
}

func (svc FilmService) GetFilm(ctx context.Context, id string) (entities.FilmEntity, error) {
	film, err := svc.repo.GetFilm(ctx, id)
	if err != nil {
		return entities.FilmEntity{}, fmt.Errorf("get film failed: %w", err)
	}
	return film, err
}

func (svc FilmService) UpdateFilm(ctx context.Context, id string, film entities.FilmEntity) error {
	err := svc.repo.UpdateFilm(ctx, id, film)
	return err
}

func (svc FilmService) DeleteFilm(ctx context.Context, id string) error {
	err := svc.repo.DeleteFilm(ctx, id)
	return err
}
