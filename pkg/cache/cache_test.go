package cache

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCache_SetAndGet(t *testing.T) {
	cache := NewCache(time.Duration(10 * time.Second))

	cache.Set("testKey", "testValue")

	cacheValue, exists := cache.Get("testKey")

	assert.True(t, exists, "cache value should exist")

	assert.Equal(t, "testValue", cacheValue)
}

func TestCache_GetMissingKey(t *testing.T) {
	cache := NewCache(time.Duration(10 * time.Second))

	cacheValue, exists := cache.Get("noExistKey")

	assert.False(t, exists, "cache value of not existed key exists")

	assert.Nil(t, cacheValue, "cache value of not existed key has value")
}

func TestCache_Expiry(t *testing.T) {
	cache := NewCache(time.Duration(100 * time.Millisecond))

	cache.Set("testKey", "testValue")
	time.Sleep(200 * time.Millisecond)

	cacheValue, exists := cache.Get("testKey")

	assert.False(t, exists, "cache value of expired key exists")

	assert.Nil(t, cacheValue, "cache value of expired key has value")
}

func TestCache_ConcurrentAccess(t *testing.T) {
	var wg sync.WaitGroup

	cache := NewCache(time.Duration(100 * time.Millisecond))

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			cache.Set("testKey", "testValue")

			_, _ = cache.Get("testKey")
		}()
	}

	wg.Wait()
}
