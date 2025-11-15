// Package plot defines and builds the plots available in Statsviz.
package plot

import (
	"runtime/debug"
	"runtime/metrics"
	"slices"
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

type registry struct {
	allnames     map[string]bool
	metrics      []string // every used metrics is added
	descriptions []description

	samples []metrics.Sample // lazily built, only with the metrics we need
}

var reg = sync.OnceValue(func() *registry {
	reg := &registry{
		allnames: make(map[string]bool),
	}

	for _, m := range metrics.All() {
		reg.allnames[m.Name] = true
	}

	return reg
})

func (r *registry) mustidx(metric string) int {
	if !r.allnames[metric] {
		panic(metric + ": unknown metric in " + goversion())
	}

	idx := slices.Index(r.metrics, metric)
	if idx == -1 {
		r.metrics = append(r.metrics, metric)
		idx = len(r.metrics) - 1
	}

	return idx
}

func (r *registry) read() []metrics.Sample {
	if r.samples == nil {
		r.samples = make([]metrics.Sample, len(r.metrics))
		for i := range r.samples {
			r.samples[i].Name = r.metrics[i]
		}
	}
	metrics.Read(r.samples)

	return r.samples
}

func (r *registry) register(desc description) {
	r.descriptions = append(r.descriptions, desc)
}

func mustidx(metric string) int {
	// TODO: adapter for refactoring: remove
	return reg().mustidx(metric)
}

func register(desc description) struct{} {
	// TODO: adapter for refactoring: remove
	reg().register(desc)
	return struct{}{}
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

func goversion() string {
	bnfo, ok := debug.ReadBuildInfo()
	if ok {
		return bnfo.GoVersion
	}

	return "<unknown version>"
}
