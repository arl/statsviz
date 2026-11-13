// Package plot defines and builds the plots available in Statsviz.
package plot

//go:generate go run -C internal/gen .

import (
	"runtime/metrics"
	"slices"
	"sync"
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
	allMetrics   map[string]bool // names of all known runtime/metrics metrics
	index        map[string]int  // index in samples of the metrics we read
	metrics      []string
	descriptions []description

	samples []metrics.Sample // lazily built, only with the metrics we need
}

// internalMetrics are read by transforms, see helpers.go.
var internalMetrics = []string{
	"/cpu/classes/total:cpu-seconds",
	"/sched/gomaxprocs:threads",
	"/gc/cycles/total:gc-cycles",
}

var reg = sync.OnceValue(func() *registry {
	reg := &registry{
		allMetrics: make(map[string]bool),
		index:      make(map[string]int),
	}
	for _, m := range metrics.All() {
		reg.allMetrics[m.Name] = true
	}
	for _, m := range internalMetrics {
		reg.add(m)
	}

	return reg
})

func (r *registry) add(metric string) {
	if _, ok := r.index[metric]; !ok {
		r.index[metric] = len(r.metrics)
		r.metrics = append(r.metrics, metric)
	}
}

func (r *registry) read() []metrics.Sample {
	if len(r.samples) != len(r.metrics) {
		r.samples = make([]metrics.Sample, len(r.metrics))
		for i := range r.samples {
			r.samples[i].Name = r.metrics[i]
		}
	}
	metrics.Read(r.samples)

	return r.samples
}

// register adds a plot to the registry, unless one of its metrics is unknown
// to the current Go version, or a plot with the same name has already been
// registered (which allows to define fallback plots for older Go versions).
func (r *registry) register(desc description) {
	for _, metric := range desc.metrics {
		if !r.allMetrics[metric] {
			return
		}
	}
	for _, metric := range desc.metrics {
		r.add(metric)
	}

	// Histograms need special handling.
	type heatmapLayoutFunc = func(samples []metrics.Sample) Heatmap
	if buildLayout, ok := desc.layout.(heatmapLayoutFunc); ok {
		desc.layout = buildLayout(r.read())
	}

	name := nameFromLayout(desc.layout)
	if slices.ContainsFunc(r.descriptions, func(d description) bool { return nameFromLayout(d.layout) == name }) {
		return
	}
	r.descriptions = append(r.descriptions, desc)
}

func register(desc description) struct{} {
	reg().register(desc)
	return struct{}{}
}

// sample returns the sample of a metric, which must have been registered.
func sample(samples []metrics.Sample, metric string) metrics.Value {
	idx, ok := reg().index[metric]
	if !ok {
		panic(metric + ": metric not registered")
	}
	return samples[idx].Value
}
