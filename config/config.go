package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env string `env:"ENV"`

	PostgresDBUsername string `env:"POSTGRES_DB_USERNAME"`
	PostgresDBPassword string `env:"POSTGRES_DB_PASSWORD"`
	PostgresDBHost     string `env:"POSTGRES_DB_HOST"`
	PostgresDBName     string `env:"POSTGRES_DB_NAME"`
	MigratePath        string `env:"MIGRATE_PATH"`

	ServerHost string `env:"SERVER_HOST"`

	AdminLogin string `env:"ADMIN_LOGIN"`
	AdminPass  string `env:"ADMIN_PASS"`

	AccessTokenExp  int    `env:"ACCESS_TOKEN_EXP"`
	RefreshTokenExp int    `env:"REFRESH_TOKEN_EXP"`
	Hs256Secret     string `env:"HS256_SECRET"`

	RedisDbHost     string `env:"REDIS_DB_HOST"`
	RedisDbPassword string `env:"REDIS_DB_PASSWORD"`
	RedisDbName     int    `env:"REDIS_DB_NAME"`
}

func New() (Config, error) {
	err := godotenv.Load("./config/.env")
	if err != nil {
		return Config{}, fmt.Errorf("load failed: %w", err)
	}

	cfg := Config{}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return Config{}, fmt.Errorf("parse failed: %w", err)
	}
	return cfg, nil
}

func (c *Config) GetPostgresUrl() string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s", c.PostgresDBUsername, c.PostgresDBPassword, c.PostgresDBHost, c.PostgresDBName)
}
