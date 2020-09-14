package pcache

import (
	"sync"
)

type cache map[string]interface{}

type cacheStripe struct {
	deletionProbeIdx int
	cache            cache
}

// PCache is safe for concurrent use, goroutine local cache.
// All operations run in amortized constant time.
// PCache does its best to cache items inside and do as little synchronization as possible
// but since it is cache, there is no guarantee that PCache won't evict your item after Store.
//
// PCache evicts random items if I goroutine local cache achieves maxSizePerGoroutine size.
// PCache cleans itself entirely from time to time.
//
// The zero PCache is invalid. Use NewPCache method to create PCache.
type PCache struct {
	maxSize               int
	deletionProbeBoundary int
	pool                  *sync.Pool
}

// NewPCache creates PCache with maxSizePerGoroutine.
func NewPCache(maxSizePerGoroutine uint) *PCache {
	return &PCache{
		maxSize:               int(maxSizePerGoroutine),
		deletionProbeBoundary: min(8, int(maxSizePerGoroutine)),
		pool: &sync.Pool{
			New: func() interface{} {
				return &cacheStripe{
					cache: make(cache),
				}
			},
		},
	}
}

// Load fetches (value, true) from cache associated with key or (nil, false) if it is not present.
func (p *PCache) Load(key string) (interface{}, bool) {
	stripe := p.pool.Get().(*cacheStripe)
	defer p.pool.Put(stripe)
	value, ok := stripe.cache[key]
	return value, ok
}

// Store stores value for a key in cache.
func (p *PCache) Store(key string, value interface{}) {
	stripe := p.pool.Get().(*cacheStripe)
	defer p.pool.Put(stripe)

	stripe.cache[key] = value
	if len(stripe.cache) <= p.maxSize {
		return
	}

	stripe.deletionProbeIdx++
	if stripe.deletionProbeIdx >= p.deletionProbeBoundary {
		stripe.deletionProbeIdx = 0
	}
	idx := 0
	for k := range stripe.cache {
		if idx == stripe.deletionProbeIdx {
			delete(stripe.cache, k)
			break
		}
		idx++
	}
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
