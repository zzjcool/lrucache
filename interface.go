package lrucache

import (
	"context"
)

type LRUCacheWithCtx[K comparable, V any] interface {
	handleWithCtx[K, V]
	Len() int
}

type LRUCache[K comparable, V any] interface {
	handle[K, V]
	Len() int
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
	Delete(ctx context.Context, k K) (err error)
}

type handle[K comparable, V any] interface {
	Get(k K) (v V, err error)
	Set(k K, v V) (err error)
	Delete(k K) (err error)
}
