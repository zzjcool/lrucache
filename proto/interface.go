package proto

import (
	"context"
	"time"
)

type LRU[K comparable, V any] interface {
	Get(k K) (v V, err error)
	Set(k K, v V)
	SetByExpire(k K, v V, expire time.Duration)
	Del(k K) (err error)
	Len() int
	Setting(capacity int, t time.Duration)
	RegisterRemoveHook(hook RemoveHook[K, V])
}

type RemoveHook[K comparable, V any] func(K, V)

type LRUCacheWithCtx[K comparable, V any] interface {
	handleWithCtx[K, V]
	Len() int

	SetByExpire(ctx context.Context, k K, v V, expire time.Duration) (err error)
}

type LRUCache[K comparable, V any] interface {
	handle[K, V]
	Len() int
	SetByExpire(k K, v V, expire time.Duration) (err error)
}

type ContextSource[K comparable, V any] interface {
	handleWithCtx[K, V]
}

type Source[K comparable, V any] interface {
	handle[K, V]
}

type handleWithCtx[K comparable, V any] interface {
	Get(ctx context.Context, k K) (v V, err error)
	Set(ctx context.Context, k K, v V) (err error)
	Del(ctx context.Context, k K) (err error)
}

type handle[K comparable, V any] interface {
	Get(k K) (v V, err error)
	Set(k K, v V) (err error)
	Del(k K) (err error)
}
