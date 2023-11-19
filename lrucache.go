package lrucache

import (
	"context"
	"sync"
	"time"

	"github.com/zzjcool/lrucache/internal/lru"
	"github.com/zzjcool/lrucache/proto"
)

func (c *core[K, V]) NewContextLRUCache(capacity int, ops ...Option[K, V]) proto.LRUCacheWithCtx[K, V] {
	lru := &cacheOption[K, V]{
		cap: capacity,
		lru: lru.NewLRU[K, V](capacity),
	}
	for _, op := range ops {
		op.apply(lru)
	}
	return lru
}

type cacheOption[K comparable, V any] struct {
	cap       int
	lru       proto.LRU[K, V]
	source    proto.ContextSource[K, V]
	downgrade bool
	keyLock   keyLock
}

type keyLock struct {
	locks sync.Map
}

func (kl *keyLock) Lock(key interface{}) {
	lock, _ := kl.locks.LoadOrStore(key, &sync.Mutex{})
	lock.(*sync.Mutex).Lock()
}

func (kl *keyLock) Unlock(key interface{}) {
	lock, ok := kl.locks.Load(key)
	if ok {
		lock.(*sync.Mutex).Unlock()
	}
}

func (c *cacheOption[K, V]) Get(ctx context.Context, k K) (v V, err error) {
	v, err = c.lru.Get(k)
	if err == nil {
		return v, nil
	}

	if c.source == nil {
		return v, err
	}

	// not found in memery
	if err == proto.NilErr {
		c.keyLock.Lock(k)
		defer c.keyLock.Unlock(k)
		v, err = c.lru.Get(k)
		if err == nil {
			return v, nil
		}
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

func (c *cacheOption[K, V]) SetByExpire(ctx context.Context, k K, v V, expire time.Duration) (err error) {
	c.lru.SetByExpire(k, v, expire)
	if c.source == nil {
		return
	}
	return c.source.Set(ctx, k, v)
}

func (c *cacheOption[K, V]) Set(ctx context.Context, k K, v V) (err error) {
	c.lru.Set(k, v)
	if c.source == nil {
		return
	}
	return c.source.Set(ctx, k, v)
}

func (c *cacheOption[K, V]) Del(ctx context.Context, k K) (err error) {
	err = c.lru.Del(k)
	if c.source == nil {
		return err
	}
	return c.source.Del(ctx, k)
}

func (c *cacheOption[K, V]) Len() int {
	return c.lru.Len()
}

func (c *cacheOption[K, V]) Clean() {
	c.lru.Clean()
}

func (c *core[K, V]) NewLRUCache(capacity int, ops ...Option[K, V]) proto.LRUCache[K, V] {
	l := c.NewContextLRUCache(capacity, ops...)
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

func (c *lruCache[K, V]) SetByExpire(k K, v V, expire time.Duration) (err error) {

	return c.lru.SetByExpire(context.Background(), k, v, expire)
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
func (c *lruCache[K, V]) Clean() {
	c.lru.Len()
}
