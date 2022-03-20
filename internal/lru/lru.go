package lru

import (
	"time"

	"github.com/zzjcool/lrucache/proto"
)

type lru[K comparable, V any] struct {
	cap        int
	hashMap    map[K]*linkNode[K, V]
	head, tail *linkNode[K, V]
	expireTime time.Duration
}

type linkNode[K comparable, V any] struct {
	Key       K
	Val       V
	ExpierAt  time.Time
	Pre, Next *linkNode[K, V]
}

func NewLRU[K comparable, V any](capacity int) proto.LRU[K, V] {
	head, tail := new(linkNode[K, V]), new(linkNode[K, V])
	head.Next, tail.Pre = tail, head
	return &lru[K, V]{
		hashMap: make(map[K]*linkNode[K, V]),
		head:    head,
		tail:    tail,
	}
}

func (l *lru[K, V]) Setting(capacity int, t time.Duration) {
	l.cap = capacity
	l.expireTime = t
}

func (l *lru[K, V]) Get(key K) (val V, err error) {
	if node, ok := l.hashMap[key]; ok {
		if l.expireTime == 0 || node.ExpierAt.After(time.Now()) {
			l.moveToHead(node)
			return node.Val, nil
		}
		// key was outdate
		return node.Val, proto.ExpiredErr
	}
	return val, proto.NilErr
}

func (l *lru[K, V]) Set(key K, value V) {
	node, ok := l.hashMap[key]
	if ok {
		node.Val = value
		if l.expireTime != 0 {
			node.ExpierAt = time.Now().Add(l.expireTime)
		}
		l.moveToHead(node)
		return
	} else {
		node = &linkNode[K, V]{
			Key:      key,
			Val:      value,
			ExpierAt: time.Now().Add(l.expireTime),
		}
		l.hashMap[key] = node
		l.addToHead(node)
		if l.Len() > l.cap {
			l.removeTail()
		}
	}
}

func (l *lru[K, V]) Del(k K) error {
	node, ok := l.hashMap[k]
	if !ok {
		return proto.NilErr
	}
	node.Pre.Next = node.Next
	node.Next.Pre = node.Pre
	delete(l.hashMap, k)
	return nil
}

func (this *lru[K, V]) moveToHead(node *linkNode[K, V]) {
	node.Pre.Next = node.Next
	node.Next.Pre = node.Pre
	this.addToHead(node)
}

func (l *lru[K, V]) addToHead(node *linkNode[K, V]) {
	node.Pre = l.head
	node.Next = l.head.Next
	node.Pre.Next = node
	node.Next.Pre = node
}

func (this *lru[K, V]) removeTail() {
	node := this.tail.Pre
	node.Pre.Next = node.Next
	node.Next.Pre = node.Pre
	delete(this.hashMap, node.Key)
}
func (this *lru[K, V]) Len() int { return len(this.hashMap) }
