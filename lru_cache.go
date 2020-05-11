package nimkv

import (
  "fmt"
  "log"
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

  l := &LRUCache {
    base: c,
    items: make(map[string]*list.Element),
    evictionList: list.New(),
  }

  if c.TickerPeriod > 0 {
    go l.ttlEvicter()
  }

  return l, nil
}

func (l *LRUCache) GetItem(key string) (*cacheItem, error) {
  // evictionList is being modified so a Writer Lock is required.
  l.base.rwLock.Lock()
  defer l.base.rwLock.Unlock()

  if item, ok := l.items[key]; ok {

    // If item has expired, return an error.
    if l.isItemExpired(item) {
      goto not_found
    }

    log.Printf("Fetching item \"%s\".", item.Value.(*cacheItem).Key)

    l.evictionList.MoveToFront(item)

    return item.Value.(*cacheItem), nil
  }
not_found:
  return nil, errors.New(fmt.Sprintf("Can't find any item with key %v.", key))
}

// Returns *cacheItems struct for all active (unexpired) items. Doesn't affect ordering of
// items in the LRU list.
func (l *LRUCache) GetAllItems() (*cacheItems) {
  l.base.rwLock.RLock()
  defer l.base.rwLock.RUnlock()

  items := make([]*cacheItem, 0, 10)

  for _, item := range l.items {

    // Ignore expired items.
    if !l.isItemExpired(item) {
      items = append(items, item.Value.(*cacheItem))
    }
  }

  return &cacheItems{ Items: items }
}

func (l *LRUCache) DeleteItem(key string) error {
  l.base.rwLock.Lock()
  defer l.base.rwLock.Unlock()

  if item, ok := l.items[key]; ok {
    log.Printf("Deleting item \"%s\".", item.Value.(*cacheItem).Key)
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

  if item, ok := l.items[key]; ok {
    log.Printf("Updating item \"%s\".", key)
    item.Value.(*cacheItem).Value = value
    item.Value.(*cacheItem).TTL = ttl
    item.Value.(*cacheItem).ExpirationTime = FriendlyTime(expirationTime)
    l.evictionList.MoveToFront(item)
  } else {
    log.Printf("Creating item \"%s\".", key)
    l.items[key] = l.evictionList.PushFront(&cacheItem{
      Key: key,
      Value: value,
      TTL: ttl,
      ExpirationTime: FriendlyTime(expirationTime),
    })
  }

  l.base.rwLock.Unlock()

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

  log.Println("Purging Cache.")

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

  itemKey := element.Value.(*cacheItem).Key

  log.Printf("[LRU-Evicter]: Evicting item \"%s\" from cache.", itemKey)

  delete(l.items, itemKey )
  l.evictionList.Remove(element)
}

func (l *LRUCache) isItemExpired(element *list.Element) bool {
  l.base.rwLock.RLock()
  defer l.base.rwLock.RUnlock()

  item := element.Value.(*cacheItem)

  if (item.TTL > 0) && time.Now().After(time.Time(item.ExpirationTime)) {
    return true
  }

  return false
}

func (l *LRUCache) ttlEvicter() {
  for range l.base.ticker {

    log.Println("[TTL-Evicter]: Trying to evict expired items...")

    for _, item := range l.items {
      if l.isItemExpired(item) {
        log.Printf("[TTL-Evicter]: Evicting item \"%s\" from cache.", item.Value.(*cacheItem).Key)
        l.evictItem(item)
      }
    }
  }
}
