package lrucache

import (
	"time"

	"github.com/zzjcool/lrucache/proto"
)

type Option[K comparable, V any] func(l *contextlruCache[K, V])

func (c *core[K, V]) WithExpireTime(t time.Duration) Option[K, V] {
	return func(l *contextlruCache[K, V]) {
		l.lru.Setting(l.cap, t)
	}
}

func (c *core[K, V]) WithSource(s proto.Source[K, V]) Option[K, V] {
	return func(l *contextlruCache[K, V]) {
		l.source = &sourceWrapCtx[K, V]{s}
	}
}

func (c *core[K, V]) WithContextSource(s proto.ContextSource[K, V]) Option[K, V] {
	return func(l *contextlruCache[K, V]) {
		l.source = s
	}
}

func (c *core[K, V]) WithDowngrade() Option[K, V] {
	return func(l *contextlruCache[K, V]) {
		l.downgrade = true
	}
}
