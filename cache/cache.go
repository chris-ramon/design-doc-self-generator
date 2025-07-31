package cache

import "github.com/golang/groupcache/lru"

type Cache struct {
	cache lru.Cache
}

func (c *Cache) Add(key string, value any) {
	c.cache.Add(lru.Key(key), value)
}

func (c *Cache) Get(key string) (any, bool) {
	return c.cache.Get(lru.Key(key))
}

func New() *Cache {
	c := &Cache{
		cache: lru.Cache{},
	}

	return c
}
