package plot

import (
	"runtime/metrics"
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

	// make creates the state (support struct) for the plot.
	make func(indices ...int) metricsGetter
}

var (
	registry []description

	metricDescs = metrics.All()
	metricIdx   map[string]int
)

func register(desc description) struct{} {
	registry = append(registry, desc)
	return struct{}{}
}

func init() {
	// We need a first set of sample in order to dimension and process the
	// heatmaps buckets.
	samples := make([]metrics.Sample, len(metricDescs))
	metricIdx = make(map[string]int)

	for i := range samples {
		samples[i].Name = metricDescs[i].Name
		metricIdx[samples[i].Name] = i
	}
	metrics.Read(samples)

	registry = append(registry)

	// TODO: all descriptions where the layout is a function of the samples need to be called here
}
