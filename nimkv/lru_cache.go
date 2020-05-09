package nimkv

import (
  "fmt"
  "container/list"
  "errors"
  "time"
)

// Simple LRU Cache.
type LRUCache struct {
  base *CacheBase
  items map[interface{}]*list.Element
  evictionList *list.List
}

func NewLRUCache(c *CacheBase) (*LRUCache, []error) {
  errorList := c.checkAndSetFields()

  if len(errorList) > 0 {
    return nil, errorList
  }

  return &LRUCache {
    base: c,
    items: make(map[interface{}]*list.Element),
    evictionList: list.New(),
  }, nil
}

func (l *LRUCache) GetItem(key interface{}) (*cacheItem, error) {
  if item, ok := l.items[key]; ok {

    // If item has expired, delete it and return an error.
    if l.isItemExpired(item) {
      l.evictItem(item)
      goto not_found
    }

    l.evictionList.MoveToFront(item)
    return item.Value.(*cacheItem), nil
    // return item.Value.(*cacheItem).value, nil
  }
not_found:
  return nil, errors.New(fmt.Sprintf("Can't find any item with key %v.", key))
}

func (l *LRUCache) DeleteItem(key interface{}) error {
  if item, ok := l.items[key]; ok {
    l.evictItem(item)
    return nil
  }

  return errors.New(fmt.Sprintf("Can't find any item with key %v.", key))
}

// Adds item to LRUCache. If item already exists, it will be overriden with new value.
// If capacity is exceeded, least recently used item will be evicted.
func (l *LRUCache) SetItemWithExpiry(key interface{}, value interface{}, ttl time.Duration) {
  expirationTime := time.Now().Add(time.Duration(ttl) * time.Second)

  if item, ok := l.items[key]; ok {
    item.Value.(*cacheItem).value = value
    l.evictionList.MoveToFront(item)
  } else {
    l.items[key] = l.evictionList.PushFront(&cacheItem{
      key: key,
      value: value,
      ttl: ttl,
      expirationTime: expirationTime,
    })
  }

  if int32(l.evictionList.Len()) > l.base.Capacity {
    l.EvictNItemsByLRU(1)
  }
}

// Adds item to LRUCache with ttl of 0.
func (l *LRUCache) SetItem(key, value interface{}) {
  l.SetItemWithExpiry(key, value, 0)
}

// Evicts n items from cache by LRU eviction policy.
func (l *LRUCache) EvictNItemsByLRU(n int) {
  for l.evictionList.Len() > 0 && n > 0 {
    // delete(l.items, l.evictionList.Back().Value.(*cacheItem).key)
    // l.evictionList.Remove(l.evictionList.Back())
    l.evictItem(l.evictionList.Back())
    n--
  }
}

// Evicts specified item from cache.
func (l *LRUCache) evictItem(item *list.Element) {
  delete(l.items, item.Value.(*cacheItem).key)
  l.evictionList.Remove(item)
}

func (l *LRUCache) isItemExpired(item *list.Element) bool {
  if time.Now().After(item.Value.(*cacheItem).expirationTime) {
    return true
  }

  return false
}
