package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/erfangho/url-shortener/internal/config"
)

type RedisCache struct {
	client *config.RedisClient
	ttl    time.Duration
}

func NewRedisCache(client *config.RedisClient, ttl time.Duration) *RedisCache {
	return &RedisCache{
		client: client,
		ttl:    ttl,
	}
}

func (c *RedisCache) Get(ctx context.Context, key string, dest any) (bool, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return false, nil
	}
	return true, json.Unmarshal([]byte(val), dest)
}

func (c *RedisCache) Set(ctx context.Context, key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, c.ttl).Err()
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
