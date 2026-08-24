package cache

import (
	"log/slog"
	"sync"
	"time"
)

type cacheEntry struct {
	value     any
	ExpiresAt time.Time
}

type Cache struct {
	items map[string]*cacheEntry
	ttl   time.Duration
	mu    sync.RWMutex
}

func NewCache(ttl time.Duration) *Cache {
	c := &Cache{
		items: make(map[string]*cacheEntry),
		ttl:   ttl,
	}

	go c.cleanup()

	return c
}

func (c *Cache) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = &cacheEntry{
		value:     value,
		ExpiresAt: time.Now().Add(c.ttl),
	}

	slog.Info(
		"url cached",
		"key", key,
		"expires_at", c.items[key].ExpiresAt,
	)
}

func (c *Cache) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cache, exists := c.items[key]

	if !exists {
		return nil, false
	}

	if cache.ExpiresAt.Before(time.Now()) {
		delete(c.items, key)
		return nil, false
	}

	slog.Info("url cache get", "key", key)

	return cache.value, true
}

func (c *Cache) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		for key, value := range c.items {
			if value.ExpiresAt.Before(time.Now()) {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}
