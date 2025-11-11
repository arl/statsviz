// Package plot defines and builds the plots available in Statsviz.
package plot

import (
	"runtime/debug"
	"runtime/metrics"
	"sync"
	"time"
)

type tag = string

const (
	tagGC        tag = "gc"
	tagScheduler tag = "scheduler"
	tagCPU       tag = "cpu"
	tagMisc      tag = "misc"
)

type description struct {
	metrics []string
	layout  any

	// getvalues creates the state (support struct) for the plot.
	getvalues func() getvalues
}

var (
	registry []description

	metricIdx map[string]int
)

var initIndices = sync.OnceValue(func() []metrics.Sample {
	all := metrics.All()

	metricIdx = make(map[string]int, len(all))

	samples := make([]metrics.Sample, len(all))
	for i := range samples {
		samples[i].Name = all[i].Name
		metricIdx[samples[i].Name] = i
	}
	metrics.Read(samples)
	return samples
})

func mustidx(metric string) int {
	_ = initIndices()
	idx, ok := metricIdx[metric]
	if !ok {
		bnfo, ok := debug.ReadBuildInfo()
		if ok {
			panic(metric + ": unknown metric in " + bnfo.GoVersion)
		}
		panic(metric + ": unknown metric in current go version")
	}
	return idx
}

func register(desc description) struct{} {
	registry = append(registry, desc)
	return struct{}{}
}

func init() {
	type heatmapLayoutFunc = func(samples []metrics.Sample) Heatmap

	samples := initIndices()
	for i := range registry {
		desc := &registry[i]
		if hm, ok := desc.layout.(heatmapLayoutFunc); ok {
			desc.layout = hm(samples)
		}
	}
}

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
