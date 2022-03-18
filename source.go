package lrucache

import "context"

type EmptyContextSource[K comparable, V any] struct {
}

func (s *EmptyContextSource[K, V]) Get(ctx context.Context, k K) (v V, err error) {
	return v, NilErr
}
func (s *EmptyContextSource[K, V]) Set(ctx context.Context, k K, v V) (err error) {
	return
}
func (s *EmptyContextSource[K, V]) Delete(ctx context.Context, k K) (err error) {
	return
}

type EmptySource[K comparable, V any] struct {
}

func (s *EmptySource[K, V]) Get(k K) (v V, err error) {
	return v, NilErr
}
func (s *EmptySource[K, V]) Set(k K, v V) (err error) {
	return
}
func (s *EmptySource[K, V]) Delete(k K) (err error) {
	return
}

type sourceWrapCtx[K comparable, V any] struct {
	Source[K, V]
}

func (s *sourceWrapCtx[K, V]) Get(ctx context.Context, k K) (v V, err error) {
	return s.Source.Get(k)
}
func (s *sourceWrapCtx[K, V]) Set(ctx context.Context, k K, v V) (err error) {
	return s.Source.Set(k, v)
}
func (s *sourceWrapCtx[K, V]) Delete(ctx context.Context, k K) (err error) {
	return s.Source.Delete(k)
}
