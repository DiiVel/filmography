package redis

import (
	"filmography/config"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type Redis struct {
	client *redis.Client
	cfg    config.Config
}

func New(cfg config.Config) (Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisDbHost,
		Password: cfg.RedisDbPassword,
		DB:       cfg.RedisDbName,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return Redis{}, fmt.Errorf("client ping failed: %w", err)
	}
	return Redis{client, cfg}, nil
}

func (r Redis) AddToken(token string, expired time.Duration) error {
	err := r.client.Set(token, true, expired).Err()
	if err != nil {
		return fmt.Errorf("client set failed: %w", err)
	}
	return nil
}

func (r Redis) GetToken(token string) bool {
	val := r.client.Get(token).Val()
	return val == ""
}

func (r Redis) Close() error {
	return r.client.Close()
}
