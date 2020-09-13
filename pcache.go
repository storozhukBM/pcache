package pcache

import (
	"sync"
)

type Key struct {
	File string
	Line uint64
}

type PCache struct {
	maxSize               int
	deletionProbeBoundary int
	pool                  *sync.Pool
}

type cacheStripe struct {
	deletionProbeIdx int
	cache            map[Key]string
}

func NewPCache(maxSize uint) *PCache {
	return &PCache{
		maxSize:               int(maxSize),
		deletionProbeBoundary: min(10, int(maxSize)),
		pool: &sync.Pool{
			New: func() interface{} {
				return &cacheStripe{
					cache: make(map[Key]string),
				}
			},
		},
	}
}

func (p *PCache) Get(k Key) (string, bool) {
	stripe := p.pool.Get().(*cacheStripe)
	defer p.pool.Put(stripe)
	result, ok := stripe.cache[k]
	return result, ok
}

func (p *PCache) Set(k Key, v string) {
	stripe := p.pool.Get().(*cacheStripe)
	defer p.pool.Put(stripe)

	if len(stripe.cache) >= p.maxSize {
		stripe.deletionProbeIdx++
		if stripe.deletionProbeIdx >= p.maxSize {
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
	stripe.cache[k] = v
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
