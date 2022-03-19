package lrucache

import (
	"context"

	"github.com/zzjcool/lrucache/proto"
)

type EmptyContextSource[K comparable, V any] struct {
}

func (s *EmptyContextSource[K, V]) Get(ctx context.Context, k K) (v V, err error) {
	return v, proto.NilErr
}
func (s *EmptyContextSource[K, V]) Set(ctx context.Context, k K, v V) (err error) {
	return
}
func (s *EmptyContextSource[K, V]) Del(ctx context.Context, k K) (err error) {
	return
}

type EmptySource[K comparable, V any] struct {
}

func (s *EmptySource[K, V]) Get(k K) (v V, err error) {
	return v, proto.NilErr
}
func (s *EmptySource[K, V]) Set(k K, v V) (err error) {
	return
}
func (s *EmptySource[K, V]) Del(k K) (err error) {
	return
}

type sourceWrapCtx[K comparable, V any] struct {
	proto.Source[K, V]
}

func (s *sourceWrapCtx[K, V]) Get(ctx context.Context, k K) (v V, err error) {
	return s.Source.Get(k)
}
func (s *sourceWrapCtx[K, V]) Set(ctx context.Context, k K, v V) (err error) {
	return s.Source.Set(k, v)
}
func (s *sourceWrapCtx[K, V]) Del(ctx context.Context, k K) (err error) {
	return s.Source.Del(k)
}
