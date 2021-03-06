package lru

import (
	"sync"
	"time"

	"github.com/zzjcool/lrucache/proto"
)

type lru[K comparable, V any] struct {
	cap               int
	hashMap           map[K]*linkNode[K, V]
	head, tail        *linkNode[K, V]
	defaultExpireTime time.Duration
	mu                sync.Mutex

	removeHook proto.RemoveHook[K, V]
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
		cap:     capacity,
	}
}

func (l *lru[K, V]) Setting(capacity int, t time.Duration) {
	l.cap = capacity
	l.defaultExpireTime = t
}

func (l *lru[K, V]) RegisterRemoveHook(hook proto.RemoveHook[K, V]) {
	l.removeHook = hook
}

func (l *lru[K, V]) Get(key K) (val V, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if node, ok := l.hashMap[key]; ok {
		if l.defaultExpireTime == 0 || node.ExpierAt.After(time.Now()) {
			l.moveToHead(node)
			return node.Val, nil
		}
		// key was outdate
		return node.Val, proto.ExpiredErr
	}
	return val, proto.NilErr
}

func (l *lru[K, V]) SetByExpire(key K, value V, expire time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()
	node, ok := l.hashMap[key]
	if ok {
		node.Val = value
		if expire != 0 {
			node.ExpierAt = time.Now().Add(expire)
		}
		l.moveToHead(node)
		return
	} else {
		node = &linkNode[K, V]{
			Key: key,
			Val: value,
		}
		if expire != 0 {
			node.ExpierAt = time.Now().Add(expire)
		}
		l.hashMap[key] = node
		l.addToHead(node)
		if l.Len() > l.cap {
			l.removeTail()
		}
	}
}

func (l *lru[K, V]) Set(key K, value V) {
	l.SetByExpire(key, value, l.defaultExpireTime)
}

func (l *lru[K, V]) Del(k K) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	node, ok := l.hashMap[k]
	if !ok {
		return proto.NilErr
	}
	node.Pre.Next = node.Next
	node.Next.Pre = node.Pre
	delete(l.hashMap, k)
	if l.removeHook != nil {
		l.removeHook(node.Key, node.Val)
	}
	return nil
}

func (l *lru[K, V]) Clean() {
	l.mu.Lock()
	defer l.mu.Unlock()
	for k, node := range l.hashMap {
		node.Pre.Next = node.Next
		node.Next.Pre = node.Pre
		delete(l.hashMap, k)
		if l.removeHook != nil {
			l.removeHook(node.Key, node.Val)
		}
	}
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

func (l *lru[K, V]) removeTail() {
	node := l.tail.Pre
	node.Pre.Next = node.Next
	node.Next.Pre = node.Pre
	delete(l.hashMap, node.Key)
	if l.removeHook != nil {
		l.removeHook(node.Key, node.Val)
	}
}
func (this *lru[K, V]) Len() int { return len(this.hashMap) }
