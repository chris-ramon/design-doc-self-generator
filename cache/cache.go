package cache

import (
	"sync"

	"github.com/golang/groupcache/lru"
)

type Cache struct {
	// mu guards cache.
	mu sync.RWMutex

	cache *lru.Cache
}

func (c *Cache) Add(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache.Add(lru.Key(key), value)
}

func (c *Cache) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cache.Get(lru.Key(key))
}

func New() *Cache {
	// Initializes LRU cache with a reasonable 1k as max entries.
	cache := lru.New(1000)

	c := &Cache{
		cache: cache,
	}

	return c
}
