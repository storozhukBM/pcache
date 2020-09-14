package pcache_test

import (
	"github.com/storozhukBM/pcache"
	"sync"
	"sync/atomic"
	"testing"
)

func TestConcurrentHitRatioTest(t *testing.T) {
	ratio := runConcurrentHitRatioTest()
	t.Logf("act ratio: %+v", ratio)
	if ratio < 0.1 {
		t.Fatalf("even with race detector enabled ration should be 0.1 was: %v", ratio)
	}
}

func runConcurrentHitRatioTest() float64 {
	hits := uint64(0)
	cache := pcache.NewPCache(4)
	keys := []string{"1", "2", "3", "4", "5"}

	sync.Map{}

	workers := 16
	iterations := 1000

	for i := 0; i < iterations; i++ {
		var wg sync.WaitGroup
		wg.Add(workers)
		for j := 0; j < workers; j++ {
			go func() {
				defer wg.Done()
				for _, key := range keys {
					cache.Store(key, true)
				}
				for _, key := range keys {
					v, ok := cache.Load(key)
					if ok && v.(bool) {
						atomic.AddUint64(&hits, 1)
					}
				}
			}()
		}
		wg.Wait()
	}

	ratio := float64(hits) / float64(workers*iterations*len(keys))
	return ratio
}
