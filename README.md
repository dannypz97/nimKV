# nimKV
nimKV is my shot at creating an extremely basic in-memory key-value store in Golang.
Inspired by projects like gcache, go-cache and golang-lru.

## Supported features: -
- Add / Retrieve / Delete Items (Items are key-value pairs).
- Cache Purging.
- A simple server that exposes RESTful endpoints for the above features.
- LRU eviction.
- TTL eviction.

A TTL can be specified when adding items to the cache. TTL eviction works by having a ticker run in the background every 'x' seconds. 'x' is specified as TickerPeriod under config.yaml.

### NOTE: -
- Some operations MIGHT NOT be thread-safe. Making operations thread-safe requires acquiring and releasing reader / writer locks correctly. Since the cache uses only one map construct to store values, acquiring a lock against the entire cache hurts performance.
A clever optimisation would be to implement some kind of sharding, and then map keys to the right shard. Then, locks could be acquired against a shard and not the entire cache.

- I have added just a few unit tests and benchmarks to get an idea of the testing functionality offered by Go. Code coverage comes to about 45%.
