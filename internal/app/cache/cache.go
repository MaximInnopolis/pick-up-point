package cache

import (
	"sync"
	"time"
)

type Cached[V any] struct {
	expiredAt time.Time
	value     V
}

func NewCached[V any](expiredAt time.Time, value V) *Cached[V] {
	return &Cached[V]{
		expiredAt: expiredAt,
		value:     value,
	}
}

func (c *Cached[V]) Expired(now time.Time) bool {
	return c.expiredAt.Before(now)
}

func (c *Cached[V]) Value() V {
	return c.value
}

type IMCache[K comparable, V any] struct {
	ttl  time.Duration
	data map[K]*Cached[V]
	lock sync.RWMutex
}

func NewIMCache[K comparable, V any](ttl time.Duration) *IMCache[K, V] {
	return &IMCache[K, V]{
		ttl:  ttl,
		data: make(map[K]*Cached[V]),
	}

}

func (c *IMCache[K, V]) Set(key K, value V, now time.Time) {
	wrapped := NewCached[V](now.Add(c.ttl), value)
	c.lock.Lock()
	c.data[key] = wrapped
	c.lock.Unlock()
}

func (c *IMCache[K, V]) Get(key K) (V, bool) {
	c.lock.RLock()
	val, ok := c.data[key]
	c.lock.RUnlock()

	if ok && !val.Expired(time.Now()) {
		return val.Value(), true
	}

	return (&Cached[V]{}).Value(), false
}

func (c *IMCache[K, V]) Delete(key K) {
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.data, key)
}

// InvalidateExpired invalidates all expired items from the cache.
func (c *IMCache[K, V]) InvalidateExpired() {
	c.lock.Lock()
	defer c.lock.Unlock()

	now := time.Now()
	for key, val := range c.data {
		if val.Expired(now) {
			delete(c.data, key)
		}
	}
}
