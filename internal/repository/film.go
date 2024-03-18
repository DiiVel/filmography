package repository

import (
	"context"
	"filmography/internal/entities"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func (r Repo) CreateFilm(ctx context.Context, film entities.FilmEntity) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(queryCtx, "INSERT INTO films (id, title, description, release_date, rating) VALUES($1, $2, $3, $4, $5)", film.ID, film.Title, film.Description, film.ReleaseDate, film.Rating)
	if row.Err() != nil {
		return fmt.Errorf("query row context order failed: %w", row.Err())
	}

	return nil
}

func (r Repo) GetFilms(ctx context.Context) ([]entities.FilmEntity, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(queryCtx, "SELECT * FROM films")
	if err != nil {
		return nil, fmt.Errorf("query context failed: %w", err)
	}
	defer rows.Close()

	films := make([]entities.FilmEntity, 0)

	for rows.Next() {
		film := entities.FilmEntity{}
		err := rows.Scan(&film.ID, &film.Title, &film.Description, &film.ReleaseDate, &film.Rating)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		films = append(films, film)
	}

	return films, nil
}

func (r Repo) GetFilm(ctx context.Context, id string) (entities.FilmEntity, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(queryCtx, "SELECT * FROM films")
	if row.Err() != nil {
		return entities.FilmEntity{}, fmt.Errorf("query context failed: %w", row.Err())
	}

	film := entities.FilmEntity{}
	err := row.Scan(&film.ID, &film.Title, &film.Description, &film.ReleaseDate, &film.Rating)
	if err != nil {
		return entities.FilmEntity{}, fmt.Errorf("scan failed: %w", err)
	}

	return film, nil
}

func (r Repo) UpdateFilm(ctx context.Context, id string, film entities.FilmEntity) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := r.db.ExecContext(queryCtx, "UPDATE films SET title = $1, descriprion = $2, release_date = $3, rating = $4 WHERE id = $5", film.Title, film.Description, film.ReleaseDate, film.Rating, id)
	if err != nil {
		return fmt.Errorf("exec context failed: %w", err)
	}

	num, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected failed: %w", err)
	}
	if num == 0 {
		return fmt.Errorf("film does not exists")
	}
	return nil
}

func (r Repo) DeleteFilm(ctx context.Context, id string) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := r.db.ExecContext(queryCtx, "DELETE FROM films WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("exec context failed: %w", err)
	}

	num, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected failed: %w", err)
	}
	if num == 0 {
		return fmt.Errorf("film does not exists")
	}
	return nil
}
