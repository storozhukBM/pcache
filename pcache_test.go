package pcache_test

import (
	"fmt"
	"github.com/storozhukBM/pcache"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"
)

func TestPCache_Get(t *testing.T) {
	hits := uint64(0)
	accesses := uint64(0)

	go func() {
		for {
			time.Sleep(1 * time.Second)
			cAccesses := atomic.LoadUint64(&accesses)
			cHits := atomic.LoadUint64(&hits)
			fmt.Printf(
				"accesses: %12d; hits: %12d; delta: %6d; ratio: %.9f; \n",
				cAccesses, cHits, cAccesses-cHits, float64(cHits)/float64(cAccesses),
			)
		}
	}()
	keys := []pcache.Key{
		{File: "k1", Line: 1},
		{File: "k2", Line: 1},
		{File: "k3", Line: 1},
		{File: "k4", Line: 1},
		{File: "k5", Line: 1},
		{File: "k6", Line: 1},
	}
	cache := pcache.NewPCache(5)
	for i := 0; i < 16; i++ {
		go func() {
			rnd := rand.New(rand.NewSource(time.Now().Unix()))
			for {
				atomic.AddUint64(&accesses, 1)
				k := keys[rnd.Intn(6)]
				v, ok := cache.Get(k)
				if !ok {
					cache.Set(k, "v")
					continue
				}
				if v != "v" {
					t.Fatalf("unexpected value: %v", v)
				}
				atomic.AddUint64(&hits, 1)
			}
		}()
	}

	time.Sleep(100 * time.Second)
}
