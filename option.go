package lrucache

import (
	"time"

	"github.com/zzjcool/lrucache/proto"
)

// Option configures
type Option[K comparable, V any] interface {
	apply(*cacheOption[K, V])
}

type funcOption[K comparable, V any] func(l *cacheOption[K, V])

func (f funcOption[K, V]) apply(l *cacheOption[K, V]) {
	f(l)
}

func (c *core[K, V]) WithExpireTime(t time.Duration) Option[K, V] {
	return funcOption[K, V](func(l *cacheOption[K, V]) {
		l.lru.Setting(l.cap, t)
	})
}

func (c *core[K, V]) WithSource(s proto.Source[K, V]) Option[K, V] {
	return funcOption[K, V](func(l *cacheOption[K, V]) {
		l.source = &sourceWrapCtx[K, V]{s}
	})
}

func (c *core[K, V]) WithContextSource(s proto.ContextSource[K, V]) Option[K, V] {
	return funcOption[K, V](func(l *cacheOption[K, V]) {
		l.source = s
	})
}

func (c *core[K, V]) WithDowngrade() Option[K, V] {
	return funcOption[K, V](func(l *cacheOption[K, V]) {
		l.downgrade = true
	})
}
