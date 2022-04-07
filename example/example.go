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
