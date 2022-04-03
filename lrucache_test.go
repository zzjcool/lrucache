package lrucache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewContextLRUCache(t *testing.T) {
	key, val := "key", "val"
	core := New[string, string]()
	l := core.NewContextLRUCache(100)
	l.Set(context.Background(), key, val)
	gotVal, err := l.Get(context.Background(), key)
	assert.Equal(t, val, gotVal)
	assert.Equal(t, nil, err)

}
func TestNewLRUCache(t *testing.T) {
	key, val := "key", "val"
	core := New[string, string]()
	l := core.NewLRUCache(100)
	l.Set(key, val)
	gotVal, err := l.Get(key)
	assert.Equal(t, val, gotVal)
	assert.Equal(t, nil, err)

}

func Use() {
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
