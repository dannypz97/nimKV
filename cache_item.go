package nimkv

import (
  "fmt"
  "time"
)

// For more human-readable time.
type FriendlyTime time.Time

func (ft FriendlyTime) MarshalJSON() ([]byte, error) {
    var timeString string
    var zeroTime time.Time

    if zeroTime == time.Time(ft) {
      timeString = ""
    } else {
      timeString = time.Time(ft).Format(time.UnixDate)
    }

    stamp := fmt.Sprintf("\"%s\"", timeString)
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
