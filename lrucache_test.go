package lrucache

import (
	"testing"
	"time"

	"google.golang.org/grpc"
)

func TestNewLRUCacheWithCtx(t *testing.T) {
	Testtx()
}

func Testtx() {
	grpc.Dial("", grpc.WithAuthority(""))
	core := New[string, string]()
	core.NewContextLRUCache(100,
		core.WithExpireTime(time.Hour),
		core.WithDowngrade(),
		core.WithContextSource(&TestSource{}),
	)
	core.NewLRUCache(100,
		core.WithExpireTime(time.Second))
}

type TestSource struct {
	EmptyContextSource[string, string]
}
