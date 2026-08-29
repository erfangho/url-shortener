package config

import (
	"context"
	"errors"
	"os"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	*redis.Client
}

func NewRedisClient() (*RedisClient, error) {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	redisClient := redis.NewClient(&redis.Options{Addr: redisAddr})

	err := redisClient.Ping(context.Background()).Err()

	if err != nil {
		return nil, errors.New("redis client error")
	}

	return &RedisClient{
		redisClient,
	}, nil
}
