package cache

import (
	"github.com/puzpuzpuz/xsync/v4"
)

// Cache is a typed cache implementation, leveraging xsync.Map for concurrency
// safety.
type Cache[K comparable, V any] struct {
	cac *xsync.Map[K, V]
}

func New[K comparable, V any]() Interface[K, V] {
	return &Cache[K, V]{
		cac: xsync.NewMap[K, V](),
	}
}

// Create stores the given key-value pair, if it does not already exist, and
// returns whether the provided key already existed before. Create uses
// xsync.Map.LoadOrStore().
func (c *Cache[K, V]) Create(key K, val V) bool {
	_, exi := c.cac.LoadOrStore(key, val)
	return exi
}

// Delete simply removes the given keys from the typed cache if at least one key
// is explicitely provided. If the variadic argument remains empty, then Delete
// purges all internal state from the underlying sync map. Delete uses a
// write-lock via xsync.Map.Clear() or xsync.Map.Delete().
func (c *Cache[K, V]) Delete(key ...K) {
	for _, x := range key {
		c.cac.Delete(x)
	}

	if len(key) == 0 {
		c.cac.Clear()
	}
}

// Exists returns whether the given key is already set. Exists uses
// xsync.Map.Load().
func (c *Cache[K, V]) Exists(key K) bool {
	_, exi := c.cac.Load(key)
	return exi
}

// Length returns the amount of key-value pairs currently maintained in the
// underlying cache. Length uses xsync.Map.Size().
func (c *Cache[K, V]) Length() int {
	return c.cac.Size()
}

// Search returns the value of the given key, whether that key exists or not.
// Search uses a read-lock via xsync.Map.Load().
func (c *Cache[K, V]) Search(key K) (V, bool) {
	return c.cac.Load(key)
}

// Update initializes the given key-value pair or overwrites it in case the
// given key existed before. Update uses xsync.Map.Store().
func (c *Cache[K, V]) Update(key K, val V) {
	c.cac.Store(key, val)
}
