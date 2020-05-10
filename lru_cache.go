package nimkv

import (
  "fmt"
  "container/list"
  "errors"
  "time"
)

// Simple LRU Cache.
type LRUCache struct {
  base *cacheBase
  items map[string]*list.Element
  evictionList *list.List
}

func NewLRUCache(c *cacheBase) (*LRUCache, []error) {
  errorList := c.checkAndSetFields()

  if len(errorList) > 0 {
    return nil, errorList
  }

  c.Type = "LRU"

  return &LRUCache {
    base: c,
    items: make(map[string]*list.Element),
    evictionList: list.New(),
  }, nil
}

func (l *LRUCache) GetItem(key string) (*cacheItem, error) {
  l.base.rwLock.RLock()
  defer l.base.rwLock.RUnlock()

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

func (l *LRUCache) DeleteItem(key string) error {
  l.base.rwLock.Lock()
  defer l.base.rwLock.Unlock()

  if item, ok := l.items[key]; ok {
    l.evictItem(item)
    return nil
  }

  return errors.New(fmt.Sprintf("Can't find any item with key %v.", key))
}

// Adds item to LRUCache. If item already exists, it will be overriden with new value.
// If capacity is exceeded, least recently used item will be evicted.
func (l *LRUCache) SetItemWithExpiry(key string, value interface{}, ttl time.Duration) {
  var expirationTime  time.Time

  if ttl > 0 {
    expirationTime = time.Now().Add(time.Duration(ttl) * time.Second)
  }

  l.base.rwLock.Lock()
  defer l.base.rwLock.Unlock()

  if item, ok := l.items[key]; ok {
    item.Value.(*cacheItem).Value = value
    item.Value.(*cacheItem).TTL = ttl
    item.Value.(*cacheItem).ExpirationTime = FriendlyTime(expirationTime)
    l.evictionList.MoveToFront(item)
  } else {
    l.items[key] = l.evictionList.PushFront(&cacheItem{
      Key: key,
      Value: value,
      TTL: ttl,
      ExpirationTime: FriendlyTime(expirationTime),
    })
  }

  if int32(l.evictionList.Len()) > l.base.Capacity {
    l.evictNItems(1)
  }
}

// Adds item to LRUCache with ttl of 0.
func (l *LRUCache) SetItem(key string, value interface{}) {
  l.SetItemWithExpiry(key, value, 0)
}

func (l *LRUCache) Purge() {
  l.base.rwLock.Lock()
  defer l.base.rwLock.Unlock()

  l.evictionList.Init()

  // No references to map after reassignment so it should be garbage collected (eventually)
  l.items = make(map[string]*list.Element)
}

// Evicts n items from cache by LRU eviction policy.
func (l *LRUCache) evictNItems(n int) {
  for l.evictionList.Len() > 0 && n > 0 {
    l.evictItem(l.evictionList.Back())
    n--
  }
}

// Evicts specified item from cache.
func (l *LRUCache) evictItem(element *list.Element) {
  l.base.rwLock.Lock()
  defer l.base.rwLock.Unlock()

  delete(l.items, element.Value.(*cacheItem).Key)
  l.evictionList.Remove(element)
}

func (l *LRUCache) isItemExpired(element *list.Element) bool {
  l.base.rwLock.Lock()
  defer l.base.rwLock.Unlock()

  item := element.Value.(*cacheItem)

  if (item.TTL > 0) && time.Now().After(time.Time(item.ExpirationTime)) {
    return true
  }

  return false
}
