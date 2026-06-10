package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type LocalCache struct {
	client *cache.Cache
}

func NewLocalCache() *LocalCache {
	// Initialize with a default TTL of 5 minutes and cleanup interval of 10 minutes.
	return &LocalCache{
		client: cache.New(5*time.Minute, 10*time.Minute),
	}
}

// Set stores a value in the cache with a specified TTL (duration)
func (c *LocalCache) Set(key string, value interface{}, ttl time.Duration) {
	c.client.Set(key, value, ttl)
}

// Get retrieves a value from the cache. Returns nil and false if not found or expired.
func (c *LocalCache) Get(key string) (interface{}, bool) {
	return c.client.Get(key)
}

// Delete removes a specific key from the cache.
func (c *LocalCache) Delete(key string) {
	c.client.Delete(key)
}

// DeleteByPrefix deletes all keys that start with a specific prefix.
func (c *LocalCache) DeleteByPrefix(prefix string) {
	for k := range c.client.Items() {
		if len(k) >= len(prefix) && k[:len(prefix)] == prefix {
			c.client.Delete(k)
		}
	}
}


