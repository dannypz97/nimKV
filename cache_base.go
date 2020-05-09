package nimkv

import (
  "errors"
  "sync"
  "time"
)

// cacheBase struct has fields that could be reused across various specialised cache implementations.
// No field to store items is present as the implementation of these will vary from cache to cache.
// For instance, LRU cache could use a doubly linked list for fast eviction, while LFU cache could
// use a minHeap.
type CacheBase struct {
  Capacity int32
  rwLock sync.RWMutex

  // Will be assigned to ttl field of cacheItem if no expiration specified for an item when it is created.
  DefaultTTL time.Duration
}

// cacheItem represents an item that will be stored in the cache.
type cacheItem struct {
  key interface{}
  value interface{}

  // Unit: seconds.
  // A ttl of 0 means the entry will never expire. It can still be evicted.
  ttl time.Duration
  expirationTime time.Time
}

// Validates receiver struct, and initializes some fields.
func (c *CacheBase) checkAndSetFields() []error {
  errorList := make([]error, 0, 2)

  if c.Capacity <= 0 {
    errorList = append(errorList, errors.New("Cache Capacity has to be > 0."))
  }
  if c.DefaultTTL < 0 {
    errorList = append(errorList, errors.New("DefaultTTL has to be >= 0;"))
  }

  if len(errorList) > 0 {
    return errorList
  }

  c.rwLock = sync.RWMutex{}
  return nil
}

func (c *cacheItem) Key() interface{} {
  return c.key
}

func (c *cacheItem) Value() interface{} {
  return c.value
}

func (c *cacheItem) TTL() time.Duration {
  return c.ttl
}

func (c *cacheItem) ExpirationTime() time.Time {
  return c.expirationTime
}
