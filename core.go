package lrucache

func New[K comparable, V any]() *core[K, V] {
	return &core[K, V]{}
}

type core[K comparable, V any] struct{}
