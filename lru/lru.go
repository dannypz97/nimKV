// Package lru implements a simple data structure with functionality for LRU eviction.

package lru

import (
  "fmt"
  "container/list"
  "errors"
)

type LRU struct {
  size int
  entries map[interface{}]*list.Element
  evictList *list.List
}

type listEntry struct {
  key interface{}
  value interface{}
}

func NewLRU(size int) (*LRU, error) {
  if size <= 0 {
    return nil, errors.New("LRU capacity should be >= 0.")
  }

  return &LRU {
    size: size,
    entries: make(map[interface{}]*list.Element),
    evictList: list.New(),
  }, nil
}

// Adds entry to LRU. If key already exists, it will override that entry with new value.
// If size is exceeded, it will evict most recently used entry.
func (l *LRU) Set(key, value interface{}) {
  if entry, ok := l.entries[key]; ok {
    entry.Value.(*listEntry).value = value
    l.evictList.MoveToFront(entry)
  } else {
    l.entries[key] = l.evictList.PushFront(&listEntry{key: key, value: value})
  }

  if l.evictList.Len() > l.size {
    l.EvictNEntries(1)
  }
}

func (l *LRU) Get (key interface{}) (interface{}, error) {
  if entry, ok := l.entries[key]; ok {
    l.evictList.MoveToFront(entry)
    return entry.Value.(*listEntry).value, nil
  } else {
    return nil, errors.New(fmt.Sprintf("Can't find any key %v.", key))
  }
}

func (l *LRU) EvictNEntries(n int) {
  for l.evictList.Len() > 0 && n > 0 {
    delete(l.entries, l.evictList.Back().Value.(*listEntry).key)
    l.evictList.Remove(l.evictList.Back())
    n--
  }
}
