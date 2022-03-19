package lrucache

import (
	"context"

	"github.com/zzjcool/lrucache/internal/lru"
	"github.com/zzjcool/lrucache/proto"
)

func (c *core[K, V]) NewContextLRUCache(capacity int, expand ...Option[K, V]) proto.LRUCacheWithCtx[K, V] {
	lru := &contextlruCache[K, V]{
		cap: capacity,
		lru: lru.NewLRU[K, V](capacity),
	}
	for _, e := range expand {
		e(lru)
	}
	return lru
}

type contextlruCache[K comparable, V any] struct {
	cap       int
	lru       proto.LRU[K, V]
	source    proto.ContextSource[K, V]
	downgrade bool
}

func (c *contextlruCache[K, V]) Get(ctx context.Context, k K) (v V, err error) {
	v, err = c.lru.Get(k)
	if err == nil {
		return v, nil
	}

	if c.source == nil {
		return v, err
	}

	// not found in memery
	if err == proto.NilErr {
		v, err = c.source.Get(ctx, k)
		if err != nil {
			return v, err
		}
		c.lru.Set(k, v)
		return v, nil
	}

	// expired
	if err == proto.ExpiredErr {
		var newV V
		newV, err = c.source.Get(ctx, k)
		if err != nil {
			// downgrades
			if c.downgrade {
				return v, nil
			}
			return newV, err
		}
		c.lru.Set(k, newV)
		return newV, nil
	}

	return v, err
}

func (c *contextlruCache[K, V]) Set(ctx context.Context, k K, v V) (err error) {
	c.lru.Set(k, v)
	if c.source == nil {
		return
	}
	return c.source.Set(ctx, k, v)
}

func (c *contextlruCache[K, V]) Del(ctx context.Context, k K) (err error) {
	err = c.lru.Del(k)
	if c.source == nil {
		return err
	}
	return c.source.Del(ctx, k)
}

func (c *contextlruCache[K, V]) Len() int {
	return c.lru.Len()
}

func (c *core[K, V]) NewLRUCache(capacity int, expand ...Option[K, V]) proto.LRUCache[K, V] {
	l := c.NewContextLRUCache(capacity, expand...)
	lru := &lruCache[K, V]{
		lru: l,
	}
	return lru
}

type lruCache[K comparable, V any] struct {
	lru proto.LRUCacheWithCtx[K, V]
}

func (c *lruCache[K, V]) Get(k K) (v V, err error) {

	return c.lru.Get(context.Background(), k)
}

func (c *lruCache[K, V]) Set(k K, v V) (err error) {

	return c.lru.Set(context.Background(), k, v)
}

func (c *lruCache[K, V]) Del(k K) (err error) {

	return c.lru.Del(context.Background(), k)
}

func (c *lruCache[K, V]) Len() int {
	return c.lru.Len()
}
