package config

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	*redis.Client
}

func NewRedisClient() (*RedisClient, error) {
	redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	err := redisClient.Ping(context.Background()).Err()

	if err != nil {
		return nil, errors.New("redis client error")
	}

	return &RedisClient{
		redisClient,
	}, nil
}
