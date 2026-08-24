package cache

import (
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
	return &Cache{
		items: make(map[string]*cacheEntry),
		ttl:   ttl,
	}
}

func (c *Cache) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = &cacheEntry{
		value:     value,
		ExpiresAt: time.Now().Add(c.ttl),
	}
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

	return cache.value, true
}
