package repository

import (
	"context"
	"filmography/internal/entities"
	"fmt"
	"time"
)

func (r Repo) CreateUser(ctx context.Context, user entities.UserEntity) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(queryCtx, "INSERT INTO users (id, username) VALUES($1, $2)", user.ID, user.Username)
	if row.Err() != nil {
		return fmt.Errorf("query row context order failed: %w", row.Err())
	}

	return nil
}

func (r Repo) GetUsers(ctx context.Context) ([]entities.UserEntity, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(queryCtx, "SELECT * FROM users")
	if err != nil {
		return nil, fmt.Errorf("query context failed: %w", err)
	}
	defer rows.Close()

	users := make([]entities.UserEntity, 0)

	for rows.Next() {
		user := entities.UserEntity{}
		err := rows.Scan(&user.ID, &user.Username)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		users = append(users, user)
	}

	return users, nil
}

func (r Repo) GetUser(ctx context.Context, id string) (entities.UserEntity, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(queryCtx, "SELECT * FROM users")
	if row.Err() != nil {
		return entities.UserEntity{}, fmt.Errorf("query context failed: %w", row.Err())
	}

	user := entities.UserEntity{}
	err := row.Scan(&user.ID, &user.Username)
	if err != nil {
		return entities.UserEntity{}, fmt.Errorf("scan failed: %w", err)
	}

	return user, nil
}

func (r Repo) UpdateUser(ctx context.Context, id string, user entities.UserEntity) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := r.db.ExecContext(queryCtx, "UPDATE users SET username = $1  WHERE id = $2", user.Username, id)
	if err != nil {
		return fmt.Errorf("exec context failed: %w", err)
	}

	num, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected failed: %w", err)
	}
	if num == 0 {
		return fmt.Errorf("user does not exists")
	}
	return nil
}

func (r Repo) DeleteUser(ctx context.Context, id string) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := r.db.ExecContext(queryCtx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("exec context failed: %w", err)
	}

	num, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected failed: %w", err)
	}
	if num == 0 {
		return fmt.Errorf("user does not exists")
	}
	return nil
}
