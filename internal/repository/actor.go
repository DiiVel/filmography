package repository

import (
	"context"
	"filmography/internal/entities"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func (r Repo) CreateActor(ctx context.Context, actor entities.ActorEntity) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(queryCtx, "INSERT INTO actors (id, name, gender, birthday) VALUES($1, $2, $3, $4)", actor.ID, actor.Name, actor.Gender, actor.Birthday)
	if row.Err() != nil {
		return fmt.Errorf("query row context order failed: %w", row.Err())
	}

	return nil
}

func (r Repo) GetActors(ctx context.Context) ([]entities.ActorEntity, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(queryCtx, "SELECT * FROM actors")
	if err != nil {
		return nil, fmt.Errorf("query context failed: %w", err)
	}
	defer rows.Close()

	actors := make([]entities.ActorEntity, 0)

	for rows.Next() {
		actor := entities.ActorEntity{}
		err := rows.Scan(&actor.ID, &actor.Name, &actor.Gender, &actor.Birthday)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		actors = append(actors, actor)
	}

	return actors, nil
}

func (r Repo) GetActor(ctx context.Context, id string) (entities.ActorEntity, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(queryCtx, "SELECT * FROM actors WHERE id = $1", id)
	if row.Err() != nil {
		return entities.ActorEntity{}, fmt.Errorf("query context failed: %w", row.Err())
	}

	actor := entities.ActorEntity{}
	err := row.Scan(&actor.ID, &actor.Name, &actor.Gender, &actor.Birthday)
	if err != nil {
		return entities.ActorEntity{}, fmt.Errorf("scan failed: %w", err)
	}

	return actor, nil
}

func (r Repo) GetFilmsByActor(ctx context.Context, actorID string) ([]entities.FilmEntity, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(queryCtx, "SELECT films.* FROM films INNER JOIN actors_films ON films.id = actors_films.film_id WHERE actors_films.actor_id = $1", actorID)
	if err != nil {
		return nil, fmt.Errorf("query context failed: %w", err)
	}
	defer rows.Close()

	films := make([]entities.FilmEntity, 0)

	for rows.Next() {
		film := entities.FilmEntity{}
		err := rows.Scan(&film.ID, &film.Title, &film.ReleaseDate)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		films = append(films, film)
	}

	return films, nil
}

func (r Repo) UpdateActor(ctx context.Context, id string, actor entities.ActorEntity) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := r.db.ExecContext(queryCtx, "UPDATE actors SET name = $1, gender = $2, birthday = $3 WHERE id = $4", actor.Name, actor.Gender, actor.Birthday, id)
	if err != nil {
		return fmt.Errorf("exec context failed: %w", err)
	}

	num, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected failed: %w", err)
	}
	if num == 0 {
		return fmt.Errorf("actor does not exists")
	}
	return nil
}

func (r Repo) DeleteActor(ctx context.Context, id string) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := r.db.ExecContext(queryCtx, "DELETE FROM actors WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("exec context failed: %w", err)
	}

	num, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected failed: %w", err)
	}
	if num == 0 {
		return fmt.Errorf("actor does not exists")
	}
	return nil
}
