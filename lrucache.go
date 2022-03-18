package lrucache

import "context"

func (c *core[K, V]) NewContextLRUCache(capacity int, expand ...Option[K, V]) LRUCacheWithCtx[K, V] {
	lru := &contextlruCache[K, V]{
		lru: newLRU[K, V](capacity),
	}
	for _, e := range expand {
		e(lru)
	}
	return lru
}

type contextlruCache[K comparable, V any] struct {
	lru       *lru[K, V]
	source    ContextSource[K, V]
	downgrade bool
}

func (c *contextlruCache[K, V]) Get(ctx context.Context, k K) (v V, err error) {
	v, err = c.lru.get(k)
	if err == nil {
		return v, nil
	}

	if c.source == nil {
		return v, err
	}

	// not found in memery
	if err == NilErr {
		v, err = c.source.Get(ctx, k)
		if err != nil {
			return v, err
		}
		c.lru.set(k, v)
		return v, nil
	}

	// expired
	if err == ExpiredErr {
		var newV V
		newV, err = c.source.Get(ctx, k)
		if err != nil {
			// downgrades
			if c.downgrade {
				return v, nil
			}
			return newV, err
		}
		c.lru.set(k, newV)
		return newV, nil
	}

	return v, err
}

func (c *contextlruCache[K, V]) Set(ctx context.Context, k K, v V) (err error) {
	c.lru.set(k, v)
	if c.source == nil {
		return
	}
	return c.source.Set(ctx, k, v)
}

func (c *contextlruCache[K, V]) Delete(ctx context.Context, k K) (err error) {
	err = c.lru.delete(k)
	if c.source == nil {
		return err
	}
	return c.source.Delete(ctx, k)
}

func (c *contextlruCache[K, V]) Len() int {
	return c.lru.len()
}

func (c *core[K, V]) NewLRUCache(capacity int, expand ...Option[K, V]) LRUCache[K, V] {
	l := c.NewContextLRUCache(capacity, expand...)
	lru := &lruCache[K, V]{
		lru: l,
	}
	return lru
}

type lruCache[K comparable, V any] struct {
	lru LRUCacheWithCtx[K, V]
}

func (c *lruCache[K, V]) Get(k K) (v V, err error) {

	return c.lru.Get(context.Background(), k)
}

func (c *lruCache[K, V]) Set(k K, v V) (err error) {

	return c.lru.Set(context.Background(), k, v)
}

func (c *lruCache[K, V]) Delete(k K) (err error) {

	return c.lru.Delete(context.Background(), k)
}

func (c *lruCache[K, V]) Len() int {
	return c.lru.Len()
}
