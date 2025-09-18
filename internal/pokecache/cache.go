package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	mu         sync.RWMutex
	entries    map[string]cacheEntry
	interval   time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		entries:  make(map[string]cacheEntry),
		interval: interval,
	}
	cache.reapLoop()
	return cache
}

func (c *Cache) Add(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       value,
	}
}

func (c *Cache) Get(key string) (data []byte, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.entries[key]
	if !ok {
		return []byte{}, false
	}
	return entry.val, true
}

// Each time the interval passes it should remove any entries that are older than the interval.
func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	go func() {
		for now := range ticker.C {
			// Remove any entries older than interval
			c.mu.Lock()
			for key, entry := range c.entries {
				age := now.Sub(entry.createdAt)
				if age > c.interval {
					delete(c.entries, key)
				}
			}
			c.mu.Unlock()
		}
	}()
	// TODO: Consider how to clean up timer, for example:
	// return func() {
	// 	ticker.Stop()
	// }
}
