package nimkv

import (
  "testing"
  "time"
  "strconv"
  "log"
  // "fmt"
  "io/ioutil"
)

func init() {
  // Stop activity logs from being printed on-screen.
  log.SetOutput(ioutil.Discard)
}

func TestNewLRUCache(t *testing.T) {
  tables := []struct {
    capacity int
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

func TestIsItemPresent(t *testing.T) {
  tables := []struct {
    key string
    value interface{}
    shouldItemBeEvicted bool
  }{
    { key: "A", value: "A" },
    { key: "B", value: "B", shouldItemBeEvicted: true },
    { key: "C", value: "C", shouldItemBeEvicted: true },
    { key: "A", value: "B" },
    { key: "D", value: "D" },
    { key: "E", value: "E" },
    { key: "F", value: "F" },
  }

  cache, _ := NewLRUCache(&cacheBase{
    Capacity: 4,
    TickerPeriod: time.Duration(0) * time.Second,
  })

  for _, table := range tables {
    cache.SetItem(table.key, table.value)
  }

  for _, table := range tables {
    if table.shouldItemBeEvicted && cache.IsItemPresent(table.key) {
      t.Errorf("Item %s should have been evicted", table.key)
    }
    if !table.shouldItemBeEvicted && !cache.IsItemPresent(table.key) {
      t.Errorf("Item %s should have been present", table.key)
    }
  }
}

func TestSetItem(t *testing.T) {
  tables := []struct {
    key string
    value interface{}
    shouldBeOverwritten bool
  }{
    { key: "A", value: "A", shouldBeOverwritten: true },
    { key: "B", value: "B" },
    { key: "C", value: "C" },
    { key: "A", value: "B" },
    { key: "D", value: "D" },
  }

  cache, _ := NewLRUCache(&cacheBase{
    Capacity: len(tables),
    TickerPeriod: time.Duration(0) * time.Second,
  })

  for _, table := range tables {
    cache.SetItem(table.key, table.value)
  }

  for _, table := range tables {
    item, _ := cache.GetItem(table.key)
    itemStringValue := item.Value

    if !table.shouldBeOverwritten && (table.value != itemStringValue) {
      t.Errorf("Expected item %s to have value %s", table.key, table.value)
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
