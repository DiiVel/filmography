package repository

import (
	"database/sql"
	"errors"
	"filmography/config"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Repo struct {
	db  *sql.DB
	cfg config.Config
}

func New(cfg config.Config) (Repo, error) {
	db, err := sql.Open("pgx", cfg.GetPostgresUrl())
	if err != nil {
		return Repo{}, fmt.Errorf("open failed: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return Repo{}, fmt.Errorf("ping failed: %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return Repo{}, fmt.Errorf("with instance failed: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(cfg.MigratePath, "postgres", driver)
	if err != nil {
		return Repo{}, fmt.Errorf("new with database instance failed: %w", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return Repo{}, fmt.Errorf("migrate up failed: %w", err)
	}

	return Repo{
		db: db,
	}, nil
}

func (r Repo) Close() error {
	return r.db.Close()
}
