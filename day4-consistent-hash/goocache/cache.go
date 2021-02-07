package goocache

import (
	"sync"

	"kentliuqiao.com/goo-cache/day2-single-node/goocache/lru"
)

type cache struct {
	mu         sync.Mutex
	cache      *lru.Cache
	cacheBytes int64
}

func (c *cache) add(key string, val ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cache == nil {
		c.cache = lru.NewCache(c.cacheBytes, nil)
	}
	c.cache.Add(key, val)
}

func (c *cache) get(key string) (val ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cache == nil {
		return
	}

	if val, ok := c.cache.Get(key); ok {
		return val.(ByteView), ok
	}

	return
}
