package lrucache

import "time"

type Option[K comparable, V any] func(l *contextlruCache[K, V])

func (c *core[K, V]) WithExpireTime(t time.Duration) Option[K, V] {
	return func(l *contextlruCache[K, V]) {
		l.lru.expireTime = t
	}
}

func (c *core[K, V]) WithSource(s Source[K, V]) Option[K, V] {
	return func(l *contextlruCache[K, V]) {
		l.source = &sourceWrapCtx[K, V]{s}
	}
}

func (c *core[K, V]) WithContextSource(s ContextSource[K, V]) Option[K, V] {
	return func(l *contextlruCache[K, V]) {
		l.source = s
	}
}

func (c *core[K, V]) WithDowngrade() Option[K, V] {
	return func(l *contextlruCache[K, V]) {
		l.downgrade = true
	}
}
