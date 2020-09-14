// +build !race

package pcache_test

import (
	"github.com/storozhukBM/pcache"
	"math"
	"reflect"
	"testing"
)

func TestConcurrentHitRatioTestNoRace(t *testing.T) {
	expectedRationWithoutRaceDetector := 0.8
	ratio := runConcurrentHitRatioTest()
	t.Logf("act ratio: %+v", ratio)
	eq(t, true, math.Abs(ratio-expectedRationWithoutRaceDetector) < 0.01)
}

func TestPCache(t *testing.T) {
	cache := pcache.NewPCache(10)

	{
		v, ok := cache.Load("k1")
		eq(t, false, ok)
		eq(t, nil, v)
	}

	{
		cache.Store("k1", "value1")
		v, ok := cache.Load("k1")
		eq(t, true, ok)
		eq(t, "value1", v)
	}
}

func eq(t *testing.T, exp interface{}, act interface{}) {
	t.Helper()
	if reflect.DeepEqual(exp, act) {
		return
	}
	t.Logf("exp: %T=`%+v`", exp, exp)
	t.Logf("act: %T=`%+v`", act, act)
	t.Fatalf("assert failed")
}
