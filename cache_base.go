package nimkv

import (
  "errors"
  "sync"
  "time"
)

type Cacher interface {
  GetItem(string) (*cacheItem, error)
  DeleteItem(string) error
  SetItemWithExpiry(string, interface{}, time.Duration)
  SetItem(string, interface{})
  Purge()
}

// CacheBase struct has fields that could be reused across various specialised cache implementations.
// No field to store items is present as the implementation of these will vary from cache to cache.
// For instance, LRU cache could use a doubly linked list for fast eviction, while LFU cache could
// use a minHeap.
// CacheBase is also used to load configuration info from config.yaml.
type cacheBase struct {
  Capacity int32 `yaml:"Capacity"`

  // Cache Type represents its eviction policy. For instance, Type could equal "LRU".
  Type string `yaml:"Type"`

  rwLock sync.RWMutex
}

// Validates receiver struct, and initializes some fields.
func (c *cacheBase) checkAndSetFields() []error {
  errorList := make([]error, 0, 2)

  if c.Capacity <= 0 {
    errorList = append(errorList, errors.New("Cache Capacity has to be > 0."))
  }

  if len(errorList) > 0 {
    return errorList
  }

  c.rwLock = sync.RWMutex{}
  return nil
}
