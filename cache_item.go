package nimkv

import (
  "fmt"
  "time"
)

// For more human-readable time.
type FriendlyTime time.Time

func (t FriendlyTime) MarshalJSON() ([]byte, error) {
    stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format(time.UnixDate))
    return []byte(stamp), nil
}

// cacheItem represents an item that will be stored in the cache.
type cacheItem struct {
  Key string `json:"key"`
  Value interface{} `json:"value"`

  // Unit: seconds.
  // A ttl of 0 means the entry will never expire. It can still be evicted.
  TTL time.Duration `json:"ttl"`
  ExpirationTime FriendlyTime `json:"ExpirationTime"`
}
