package lrucache

import (
	"time"

	"github.com/zzjcool/lrucache/proto"
)

func New[K comparable, V any]() Core[K, V] {
	return &core[K, V]{}
}

type core[K comparable, V any] struct{}

type Core[K comparable, V any] interface {
	NewContextLRUCache(capacity int, ops ...Option[K, V]) proto.LRUCacheWithCtx[K, V]
	NewLRUCache(capacity int, ops ...Option[K, V]) proto.LRUCache[K, V]
	WithExpireTime(t time.Duration) Option[K, V]
	WithSource(s proto.Source[K, V]) Option[K, V]
	WithContextSource(s proto.ContextSource[K, V]) Option[K, V]
	WithDowngrade() Option[K, V]
}
