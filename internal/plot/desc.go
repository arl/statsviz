// Package plot defines and builds the plots available in Statsviz.
package plot

import (
	"math"
	"runtime/debug"
	"runtime/metrics"
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
	name    string
	tags    []tag
	metrics []string
	layout  any

	// getvalues creates the state (support struct) for the plot.
	getvalues func() getvalues
}

var (
	registry []description

	metricIdx map[string]int
)

func mustidx(metric string) int {
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
	// We need a first set of sample in order to dimension and process the
	// heatmaps buckets.
	all := metrics.All()
	samples := make([]metrics.Sample, len(all))
	metricIdx = make(map[string]int)

	for i := range samples {
		samples[i].Name = all[i].Name
		metricIdx[samples[i].Name] = i
	}
	metrics.Read(samples)

	type heatmapLayoutFunc = func(samples []metrics.Sample) Heatmap

	for i := range registry {
		desc := &registry[i]
		if hm, ok := desc.layout.(heatmapLayoutFunc); ok {
			desc.layout = hm(samples)
		}
	}
}

func deltaUint64(midx int) func([]metrics.Sample) uint64 {
	var last uint64 = math.MaxUint64
	return func(samples []metrics.Sample) uint64 {
		cur := samples[midx].Value.Uint64()
		delta := cur - last
		if last == math.MaxUint64 {
			delta = 0
		}
		last = cur
		return delta
	}
}

func ratefloat64(midx int) func(time.Time, []metrics.Sample) float64 {
	var last float64
	var lastTime time.Time

	return func(now time.Time, samples []metrics.Sample) float64 {
		cur := samples[midx].Value.Float64()

		if lastTime.IsZero() {
			last = cur
			lastTime = now
			return 0
		}

		t := now.Sub(lastTime).Seconds()
		rate := (cur - last) / t

		last = cur
		lastTime = now

		return rate
	}
}
