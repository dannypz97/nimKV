package nimkv

import (
  "testing"
  "time"
  "strconv"
  "os"
  // "log"
)

func TestNewLRUCache(t *testing.T) {
  tables := []struct {
    capacity int32
    tickerPeriod time.Duration
    errorsExpected bool
  }{
    { capacity: 0, tickerPeriod: -5, errorsExpected: true},
    { capacity: 0, tickerPeriod: 0, errorsExpected: true},
    { capacity: 0, tickerPeriod: 20, errorsExpected: true},
    { capacity: 5, tickerPeriod: 20, errorsExpected: true},
    { capacity: 5, tickerPeriod: 30, errorsExpected: false},
    { capacity: 50, tickerPeriod: 100, errorsExpected: false},
  }

  for _, table := range tables {
    cache, errs := NewLRUCache(&cacheBase{
      Capacity: table.capacity,
      TickerPeriod: table.tickerPeriod,
    })

    if (len(errs) > 0) != table.errorsExpected {
      t.Errorf("Errors weren't returned from NewLRUCache despite invalid values for cacheBase fields.")
    }

    if !table.errorsExpected && (cache == nil) {
      t.Errorf("Expected *LRUCache.")
    }
  }
}

func BenchmarkSetUniqueItemsUninterruptedForSmallCache(b *testing.B) {
  cache, _ := NewLRUCache(&cacheBase{
    Capacity: 100,
    TickerPeriod: time.Duration(0) * time.Second,
  })

  for i := 0; i < b.N; i++ {
    cache.SetItem(strconv.Itoa(i), i)
  }
}

func BenchmarkSetUniqueItemsUninterruptedForLargeCache(b *testing.B) {
  cache, _ := NewLRUCache(&cacheBase{
    Capacity: 100000,
    TickerPeriod: time.Duration(0) * time.Second,
  })

  for i := 0; i < b.N; i++ {
    cache.SetItem(strconv.Itoa(i), i)
  }
}
