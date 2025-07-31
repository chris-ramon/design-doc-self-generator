package cache

import "github.com/golang/groupcache/lru"

type Cache struct {
	cache *lru.Cache
}

func (c *Cache) Add(key string, value any) {
	c.cache.Add(lru.Key(key), value)
}

func (c *Cache) Get(key string) (any, bool) {
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
