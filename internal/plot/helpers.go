package plot

import (
	"runtime/debug"
	"time"
)

// delta returns a function that computes the delta between successive calls.
func delta[T uint64 | float64]() func(T) T {
	first := true
	var last T
	return func(cur T) T {
		delta := cur - last
		if first {
			delta = 0
			first = false
		}
		last = cur
		return delta
	}
}

// rate returns a function that computes the rate of change per second.
func rate[T uint64 | float64]() func(time.Time, T) float64 {
	var last T
	var lastTime time.Time

	return func(now time.Time, cur T) float64 {
		if lastTime.IsZero() {
			last = cur
			lastTime = now
			return 0
		}

		t := now.Sub(lastTime).Seconds()
		rate := float64(cur-last) / t

		last = cur
		lastTime = now

		return rate
	}
}

func goversion() string {
	bnfo, ok := debug.ReadBuildInfo()
	if ok {
		return bnfo.GoVersion
	}

	return "<unknown version>"
}

// cpuRate returns a function that converts a cumulative /cpu/classes metric
// into a rate, in CPU-seconds per second.
//
// The runtime only updates /cpu/classes metrics once per GC cycle, so the rate
// can't be computed between 2 successive samples: it would be 0 most of the
// time, with a spike after each GC cycle. Instead, the rate is averaged over
// the period between 2 updates, and held until the next one. Updates are
// detected, and time is measured, with /cpu/classes/total:cpu-seconds, which
// grows by GOMAXPROCS CPU-seconds per second. Before the second update, the
// returned rate is the average since the program started.
func cpuRate() func(total float64, maxprocs uint64, cur float64) float64 {
	var lastTotal, last, rate float64

	return func(total float64, maxprocs uint64, cur float64) float64 {
		if total != lastTotal {
			rate = (cur - last) / (total - lastTotal) * float64(maxprocs)
			lastTotal = total
			last = cur
		}
		return rate
	}
}

// gcRate returns a function that converts a cumulative metric into a rate of
// change per second, averaged over GC cycles.
//
// It's meant for metrics which the runtime mostly updates in bursts around GC
// cycles, such as heap allocations (counted when allocation caches are
// flushed) and frees (counted when memory is swept). The rate is computed each
// time cycles, the count of completed GC cycles, changes, and held until the
// next change.
func gcRate[T uint64 | float64]() func(now time.Time, cycles uint64, cur T) float64 {
	var (
		last       T
		lastCycles uint64
		lastTime   time.Time
		rate       float64
	)

	return func(now time.Time, cycles uint64, cur T) float64 {
		switch {
		case lastTime.IsZero():
			last, lastCycles, lastTime = cur, cycles, now
		case cycles != lastCycles:
			rate = float64(cur-last) / now.Sub(lastTime).Seconds()
			last, lastCycles, lastTime = cur, cycles, now
		}
		return rate
	}
}
