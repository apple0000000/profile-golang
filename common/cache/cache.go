package cache

import (
	"sync"
	"time"
)

type CacheItem struct {
	Value      string
	LastAccess time.Time
}

type SimpleCache struct {
	items   map[string]*CacheItem
	mutex   sync.RWMutex
	maxSize int
}

func NewSimpleCache(maxSize int) *SimpleCache {
	return &SimpleCache{
		items:   make(map[string]*CacheItem),
		maxSize: maxSize,
	}
}

func (c *SimpleCache) Get(key string) (string, bool) {
	c.mutex.RLock()
	item, exists := c.items[key]
	c.mutex.RUnlock()

	if exists {
		c.mutex.Lock()
		item.LastAccess = time.Now()
		c.mutex.Unlock()
		return item.Value, true
	}
	return "", false
}

func (c *SimpleCache) Set(key, value string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if len(c.items) >= c.maxSize && c.items[key] == nil {
		for k := range c.items {
			delete(c.items, k)
			break
		}
	}

	c.items[key] = &CacheItem{
		Value:      value,
		LastAccess: time.Now(),
	}
}

func (c *SimpleCache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.items, key)
}

func (c *SimpleCache) Cleanup(threshold time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now()
	for key, item := range c.items {
		if now.Sub(item.LastAccess) > threshold {
			delete(c.items, key)
		}
	}
}
