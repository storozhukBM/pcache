package pcache_test

import (
	"fmt"
	"hash/maphash"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/storozhukBM/pcache"
)

type cache interface {
	Load(key string) (interface{}, bool)
	Store(key string, value interface{})
}

var workerSets = []int{1, 2, 4, 8, 12, 16, 22, 28, 32}

func BenchmarkPCacheLoadStore(b *testing.B) {
	for _, r := range workerSets {
		cache := pcache.NewPCache(80)
		b.Run(fmt.Sprintf("w-%v", r), func(b *testing.B) {
			benchLoadStore(b, r, cache)
		})
	}
}

func BenchmarkMutexCacheLoadStore(b *testing.B) {
	for _, r := range workerSets {
		cache := newMutexCache(80)
		b.Run(fmt.Sprintf("w-%v", r), func(b *testing.B) {
			benchLoadStore(b, r, cache)
		})
	}
}

func BenchmarkStripedCacheLoadStore(b *testing.B) {
	for _, r := range workerSets {
		cache := newStripedMapCache(80)
		b.Run(fmt.Sprintf("w-%v", r), func(b *testing.B) {
			benchLoadStore(b, r, cache)
		})
	}
}

func benchLoadStore(b *testing.B, workers int, cache cache) {
	keys := make([]string, 0, 100)
	for i := 0; i < 100; i++ {
		keys = append(keys, strconv.Itoa(i))
	}

	count := uint64(0)
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			idx := 0
			hits := uint64(0)
			for i := 0; i < b.N; i++ {
				idx++
				if idx == len(keys) {
					idx = 0
				}
				cache.Store(keys[len(keys)-1-idx], keys[idx])
				_, ok := cache.Load(keys[idx])
				if ok {
					hits++
				}
			}
			atomic.AddUint64(&count, hits)
			wg.Done()
		}()
	}
	wg.Wait()
	if rand.Float32() < 0.00001 {
		b.Logf("hits: %v", atomic.LoadUint64(&count))
	}
}

type stripe struct {
	m     sync.RWMutex
	store map[string]interface{}
	_     [32]byte
}

type stripedMapCache struct {
	sizePerGoroutine int
	stripes          []*stripe
}

func newStripedMapCache(sizePerGoroutine uint) *stripedMapCache {
	cache := &stripedMapCache{
		sizePerGoroutine: int(sizePerGoroutine),
		stripes:          make([]*stripe, 0, 64),
	}
	for i := 0; i < 64; i++ {
		cache.stripes = append(cache.stripes, &stripe{
			m:     sync.RWMutex{},
			store: make(map[string]interface{}),
		})
	}
	return cache
}

func (m *stripedMapCache) Load(key string) (interface{}, bool) {
	var h maphash.Hash
	_, _ = h.WriteString(key)
	idx := uint64(64-1) & h.Sum64()
	stripe := m.stripes[idx]
	stripe.m.RLock()
	defer stripe.m.RUnlock()
	val, ok := stripe.store[key]
	return val, ok
}

func (m *stripedMapCache) Store(key string, value interface{}) {
	var h maphash.Hash
	_, _ = h.WriteString(key)
	stripeIdx := uint64(64-1) & h.Sum64()
	stripe := m.stripes[stripeIdx]
	stripe.m.Lock()
	defer stripe.m.Unlock()

	stripe.store[key] = value
	if len(stripe.store) <= m.sizePerGoroutine {
		return
	}
	for k := range stripe.store {
		delete(stripe.store, k)
		break
	}
}

type mutexMapCache struct {
	m                sync.RWMutex
	sizePerGoroutine int
	store            map[string]interface{}
}

func newMutexCache(sizePerGoroutine uint) *mutexMapCache {
	return &mutexMapCache{
		m:                sync.RWMutex{},
		sizePerGoroutine: int(sizePerGoroutine),
		store:            make(map[string]interface{}),
	}
}

func (m *mutexMapCache) Load(key string) (interface{}, bool) {
	m.m.RLock()
	defer m.m.RUnlock()
	val, ok := m.store[key]
	return val, ok
}

func (m *mutexMapCache) Store(key string, value interface{}) {
	m.m.Lock()
	defer m.m.Unlock()
	m.store[key] = value
	if len(m.store) <= m.sizePerGoroutine {
		return
	}
	for k := range m.store {
		delete(m.store, k)
		break
	}
}
