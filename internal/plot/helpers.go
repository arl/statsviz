package plot

import (
	"math"
	"runtime/debug"
	"runtime/metrics"
	"time"
)

// value returns the value of a metric, as a float64.
func value(samples []metrics.Sample, metric string) float64 {
	v := sample(samples, metric)
	if v.Kind() == metrics.KindUint64 {
		return float64(v.Uint64())
	}
	return v.Float64()
}

// A transform converts the successive values of a metric (or of an
// expression of metrics) into the values to plot. Each series gets its own
// transform instance, created by one of the functions below.
type transform = func(now time.Time, samples []metrics.Sample, cur float64) float64

// delta computes the difference with the previous value.
func delta() transform {
	first := true
	var last float64

	return func(_ time.Time, _ []metrics.Sample, cur float64) float64 {
		d := cur - last
		if first {
			d = 0
			first = false
		}
		last = cur
		return d
	}
}

// rate computes the rate of change per second.
func rate() transform {
	var last float64
	var lastTime time.Time

	return func(now time.Time, _ []metrics.Sample, cur float64) float64 {
		if lastTime.IsZero() {
			last, lastTime = cur, now
			return 0
		}

		r := (cur - last) / now.Sub(lastTime).Seconds()
		last, lastTime = cur, now
		return r
	}
}

// cpuRate converts a cumulative /cpu/classes metric into a rate, in
// CPU-seconds per second.
//
// The runtime only updates /cpu/classes metrics once per GC cycle, so the rate
// can't be computed between 2 successive samples: it would be 0 most of the
// time, with a spike after each GC cycle. Instead, the rate is averaged over
// the period between 2 updates, and held until the next one. Updates are
// detected, and time is measured, with /cpu/classes/total:cpu-seconds, which
// grows by GOMAXPROCS CPU-seconds per second. Before the second update, the
// returned rate is the average since the program started.
func cpuRate() transform {
	r := heldRate()
	r(0, 0, 0) // so that the first rate is the average since program start
	return func(_ time.Time, samples []metrics.Sample, cur float64) float64 {
		total := value(samples, "/cpu/classes/total:cpu-seconds")
		maxprocs := value(samples, "/sched/gomaxprocs:threads")
		// total/maxprocs is a clock, in seconds.
		return r(total, total/maxprocs, cur)
	}
}

// gcRate converts a cumulative metric into a rate of change per second,
// averaged over GC cycles.
//
// It's meant for metrics which the runtime mostly updates in bursts around GC
// cycles, such as heap allocations (counted when allocation caches are
// flushed) and frees (counted when memory is swept). The rate is computed each
// time the count of completed GC cycles changes, and held until the next
// change.
func gcRate() transform {
	r := heldRate()
	var start time.Time
	return func(now time.Time, samples []metrics.Sample, cur float64) float64 {
		if start.IsZero() {
			start = now
		}
		cycles := value(samples, "/gc/cycles/total:gc-cycles")
		return r(cycles, now.Sub(start).Seconds(), cur)
	}
}

// heldRate returns a function computing the rate of change per second of cur,
// only when version changes, and returning the last computed rate otherwise.
// clock is the current time, in seconds.
func heldRate() func(version, clock, cur float64) float64 {
	first := true
	var lastVersion, lastClock, last, r float64

	return func(version, clock, cur float64) float64 {
		if first {
			first = false
			lastVersion, lastClock, last = version, clock, cur
			return 0
		}
		if version != lastVersion && clock != lastClock {
			r = (cur - last) / (clock - lastClock)
			lastVersion, lastClock, last = version, clock, cur
		}
		return r
	}
}

// noLimit replaces math.MaxInt64, which the runtime uses to represent 'no
// limit', with 0.
func noLimit() transform {
	return func(_ time.Time, _ []metrics.Sample, cur float64) float64 {
		if cur == math.MaxInt64 {
			return 0
		}
		return cur
	}
}

func goversion() string {
	bnfo, ok := debug.ReadBuildInfo()
	if ok {
		return bnfo.GoVersion
	}

	return "<unknown version>"
}
