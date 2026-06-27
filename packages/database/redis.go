package database

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	URL      string
	Password string
	DB       int
	PoolSize int
}

func DefaultRedisConfig(url string) RedisConfig {
	return RedisConfig{
		URL:      url,
		Password: "",
		DB:       0,
		PoolSize: 10,
	}
}

func NewRedisClient(cfg RedisConfig) (*redis.Client, error) {
	opts, err := redis.ParseURL(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis URL: %w", err)
	}

	opts.Password = cfg.Password
	opts.DB = cfg.DB
	if cfg.PoolSize > 0 {
		opts.PoolSize = cfg.PoolSize
	}

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return client, nil
}

func MustRedisClient(cfg RedisConfig) *redis.Client {
	client, err := NewRedisClient(cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to create redis client: %v", err))
	}
	return client
}

func RedisHealthCheck(ctx context.Context, client *redis.Client) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return client.Ping(ctx).Err()
}

func CloseRedis(client *redis.Client) {
	if client != nil {
		client.Close()
	}
}
