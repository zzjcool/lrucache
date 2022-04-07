package lru

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/zzjcool/lrucache/proto"
)

type kv struct {
	k string
	v *val
}

var nilVal *val

type val struct {
	val string
}

func NewTestLRU(cap int, t time.Duration) proto.LRU[string, *val] {
	l := NewLRU[string, *val](cap)
	l.Setting(cap, t)
	return l
}

// TestBasicLRU Test the basic functions of LRU
func TestBasicLRU(t *testing.T) {

	tests := []struct {
		name string
		kv   kv
	}{
		{
			name: "case1",
			kv: kv{
				k: "key",
				v: &val{"value"},
			},
		},
		{
			name: "case2",
			kv: kv{
				k: "key2",
				v: &val{"value2"},
			},
		},
	}
	l := NewTestLRU(100, 0)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l.Set(tt.kv.k, tt.kv.v)
			v, err := l.Get(tt.kv.k)
			assert.Equal(t, nil, err)
			assert.Equal(t, tt.kv.v, v)
			assert.Equal(t, 1, l.Len())
			err = l.Del(tt.kv.k)
			assert.Equal(t, nil, err)
			err = l.Del(tt.kv.k)
			assert.Equal(t, proto.NilErr, err)
			v, err = l.Get(tt.kv.k)
			assert.Equal(t, nilVal, v)
			assert.Equal(t, proto.NilErr, err)
			assert.Equal(t, 0, l.Len())
		})
	}

}

// TestValueUpdate Test values are updated
func TestValueUpdate(t *testing.T) {
	l := NewTestLRU(100, 0)
	kv1 := kv{
		k: "key",
		v: &val{"value"},
	}
	kv2 := kv{
		k: "key",
		v: &val{"new value"},
	}

	l.Set(kv1.k, kv1.v)
	l.Set(kv2.k, kv2.v)

	v, err := l.Get(kv2.k)
	assert.Equal(t, kv2.v, v)
	assert.Equal(t, nil, err)
}

// TestOutdated Testing outdated key
func TestOutdated(t *testing.T) {
	cap := 100
	l := NewTestLRU(cap, 0)
	kvs := genNRandomKV(150)
	for i := 0; i < 150; i++ {
		l.Set(kvs[i].k, kvs[i].v)
	}
	for i := 0; i < 50; i++ {
		v, err := l.Get(kvs[i].k)
		assert.Equal(t, proto.NilErr, err)
		assert.Equal(t, nilVal, v)
	}
	for i := 0; i < 100; i++ {
		j := 50 + i
		v, err := l.Get(kvs[j].k)
		assert.Equal(t, nil, err)
		assert.Equal(t, kvs[j].v, v)
	}
	assert.Equal(t, cap, l.Len())
	lr, ok := l.(*lru[string, *val])
	assert.Equal(t, true, ok)
	node := lr.head
	for i := 0; i < cap; i++ {
		node = node.Next
	}
	assert.Equal(t, lr.tail, node.Next)
}

// TestExpired Test data expired
func TestExpired(t *testing.T) {
	l := NewTestLRU(100, time.Second)
	kv := kv{
		k: "key",
		v: &val{"value"},
	}
	l.Set(kv.k, kv.v)
	time.Sleep(time.Second * 2)
	v, err := l.Get(kv.k)
	assert.Equal(t, proto.ExpiredErr, err)
	assert.Equal(t, kv.v, v)

	l.Set(kv.k, kv.v)
	v, err = l.Get(kv.k)
	assert.Equal(t, nil, err)
	assert.Equal(t, kv.v, v)
}

func genNRandomKV(n int) []kv {
	ret := []kv{}
	for i := 0; i < n; i++ {
		ret = append(ret, kv{
			k: fmt.Sprint(uuid.New().String()),
			v: &val{fmt.Sprint(uuid.New().String())},
		})
	}
	return ret
}

func benchmarkSetGen(cap int, b *testing.B) {
	l := NewTestLRU(cap, 0)
	kvs := genNRandomKV(b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Set(kvs[i].k, kvs[i].v)
	}
}

func BenchmarkSet(b *testing.B) {
	benchmarkSetGen(100000000, b)
}

func benchmarkGetGen(cap int, b *testing.B) {
	l := NewTestLRU(cap, 0)
	kvs := genNRandomKV(b.N)

	for i := 0; i < b.N; i++ {
		l.Set(kvs[i].k, kvs[i].v)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Get(kvs[i].k)
	}
}
func BenchmarkGet(b *testing.B) {
	benchmarkGetGen(100000000, b)

}

// TODO TestRace thread safety test
func TestRace(t *testing.T) {
	l := NewLRU[int, int](10)
	for i := 0; i < 10000; i++ {
		go func(i int) {
			l.Set(i, 0)
		}(i)
		go func(i int) {
			l.Get(i)
		}(i)
	}
}
