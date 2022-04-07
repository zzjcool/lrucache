# Go LRU Cache

`lrucache` is an in-memory kv cache.

This is an LRU cache implemented based on Go `generics`

## Feature

* Set cache size
* Set expiration time
* No assertion required
* Custom data source
* Support for circuit breakers and downgrades available
* Support Context
* Thread safety

## Usage Example

```go
package main

import (
 "context"
 "fmt"
 "time"

 "github.com/zzjcool/lrucache"
)

func main() {
 // create a <string, string> cache core
 core := lrucache.New[string, string]()
 // simple cache
 cache := core.NewLRUCache(100, core.WithExpireTime(time.Second*10))

 cache.Set("key", "val")
 val, err := cache.Get("key")
 if err != nil {
  panic(err)
 }

 fmt.Println(val)

 // context cache
 ctxCache := core.NewContextLRUCache(100,
  core.WithExpireTime(time.Second*10),
  core.WithContextSource(&CustomSource{}),
 )

 // get data from source
 val, err = ctxCache.Get(context.Background(), "sourceKey")
 if err != nil {
  panic(err)
 }

 fmt.Println(val)
}

type CustomSource struct {
 lrucache.EmptyContextSource[string, string]
}

func (s *CustomSource) Get(ctx context.Context, key string) (string, error) {
 if key == "sourceKey" {
  return "source value", nil
 }
 return "", fmt.Errorf("not found")
}

```

## Option

* WithExpireTime(t time.Duration): Set default expire time of cache
* WithSource(s proto.Source[K, V]): Get from source when cache doesn't get data
* WithContextSource(s proto.ContextSource[K, V]): source with context
* WithDowngrade() Option[K, V]: When the data in the cache expires, continue to use the expired data if the attempt to obtain the source update data fails
