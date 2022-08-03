package statsviz

import (
	"math"
	"runtime/debug"
	"runtime/metrics"
	"sync"

	"github.com/arl/statsviz/internal/plot"
)

var allDesc []metrics.Description
var samples []metrics.Sample

var (
	metricGCCycleAuto                       int // "/gc/cycles/automatic:gc-cycles"
	metricGCCycleForced                     int // "/gc/cycles/forced:gc-cycles"
	metricGCCycleTotal                      int // "/gc/cycles/total:gc-cycles"
	metricGCHeapAllocsBySize                int // "/gc/heap/allocs-by-size:bytes"
	metricGCHeapAllocsBytes                 int // "/gc/heap/allocs:bytes",
	metricGCHeapAllocsObjects               int // "/gc/heap/allocs:objects"
	metricGCHeapFreesBySize                 int // "/gc/heap/frees-by-size:bytes"
	metricGCHeapFreesBytes                  int // "/gc/heap/frees:bytes"
	metricGCHeapFreesObjects                int // "/gc/heap/frees:objects"
	metricsGCHeapGoalBytes                  int // "/gc/heap/goal:bytes"
	metricsGCHeapObjects                    int // "/gc/heap/objects:objects"
	metricsGCTinyAllocsObjects              int // "/gc/heap/tiny/allocs:objects"
	metricsGCPauses                         int // "/gc/pauses:seconds"
	metricsMemoryClassesHeapFreeBytes       int // "/memory/classes/heap/free:bytes"
	metricsMemoryClassesHeapObjects         int // "/memory/classes/heap/objects:bytes"
	metricsMemoryClassesHeapReleasedBytes   int // "/memory/classes/heap/released:bytes"
	metricsMemoryClassesHeapStackBytes      int // "/memory/classes/heap/stacks:bytes"
	metricsMemoryClassesHeapUnusedBytes     int // "/memory/classes/heap/unused:bytes"
	metricsMemoryClassesMetadataMCacheFree  int // "/memory/classes/metadata/mcache/free:bytes"
	metricsMemoryClassesMetadataMCacheInUse int // "/memory/classes/metadata/mcache/inuse:bytes"
	metricsMemoryClassesMetadataMSpanFree   int // "/memory/classes/metadata/mspan/free:bytes"
	metricsMemoryClassesMetadataMSpanInUse  int // "/memory/classes/metadata/mspan/inuse:bytes"
	metricsMemoryClassesMetadataOther       int // "/memory/classes/metadata/other:bytes"
	metricsMemoryClassesOsStacks            int // "/memory/classes/os-stacks:bytes"
	metricsMemoryClassesOther               int // "/memory/classes/other:bytes"
	metricsMemoryClassesProfilingBuckets    int // "/memory/classes/profiling/buckets:bytes"
	metricsMemoryClassesTotal               int // "/memory/classes/total:bytes"
	metricsSchedGoroutines                  int // "/sched/goroutines:goroutines"
	metricsSchedLatencies                   int // "/sched/latencies:seconds"
)

func populateSamples() {
	// Get descriptions for all supported metrics.
	allDesc = metrics.All()

	// Fill the slice of samples to pass to metrics.Read, and fill the global
	// variables containing indices of specific metrics in the slice of sampled
	// metrics.
	samples = make([]metrics.Sample, len(allDesc))
	samplesIdx := make(map[string]int)
	for i, metric := range allDesc {
		samplesIdx[metric.Name] = i
		samples[i].Name = allDesc[i].Name
	}

	metricGCCycleAuto = samplesIdx["/gc/cycles/automatic:gc-cycles"]
	metricGCCycleForced = samplesIdx["/gc/cycles/forced:gc-cycles"]
	metricGCCycleTotal = samplesIdx["/gc/cycles/total:gc-cycles"]
	metricGCHeapAllocsBySize = samplesIdx["/gc/heap/allocs-by-size:bytes"]
	metricGCHeapAllocsBytes = samplesIdx["/gc/heap/allocs:bytes"]
	metricGCHeapAllocsObjects = samplesIdx["/gc/heap/allocs:objects"]
	metricGCHeapFreesBySize = samplesIdx["/gc/heap/frees-by-size:bytes"]
	metricGCHeapFreesBytes = samplesIdx["/gc/heap/frees:bytes"]
	metricGCHeapFreesObjects = samplesIdx["/gc/heap/frees:objects"]
	metricsGCHeapGoalBytes = samplesIdx["/gc/heap/goal:bytes"]
	metricsGCHeapObjects = samplesIdx["/gc/heap/objects:objects"]
	metricsGCTinyAllocsObjects = samplesIdx["/gc/heap/tiny/allocs:objects"]
	metricsGCPauses = samplesIdx["/gc/pauses:seconds"]
	metricsMemoryClassesHeapFreeBytes = samplesIdx["/memory/classes/heap/free:bytes"]
	metricsMemoryClassesHeapObjects = samplesIdx["/memory/classes/heap/objects:bytes"]
	metricsMemoryClassesHeapReleasedBytes = samplesIdx["/memory/classes/heap/released:bytes"]
	metricsMemoryClassesHeapStackBytes = samplesIdx["/memory/classes/heap/stacks:bytes"]
	metricsMemoryClassesHeapUnusedBytes = samplesIdx["/memory/classes/heap/unused:bytes"]
	metricsMemoryClassesMetadataMCacheFree = samplesIdx["/memory/classes/metadata/mcache/free:bytes"]
	metricsMemoryClassesMetadataMCacheInUse = samplesIdx["/memory/classes/metadata/mcache/inuse:bytes"]
	metricsMemoryClassesMetadataMSpanFree = samplesIdx["/memory/classes/metadata/mspan/free:bytes"]
	metricsMemoryClassesMetadataMSpanInUse = samplesIdx["/memory/classes/metadata/mspan/inuse:bytes"]
	metricsMemoryClassesMetadataOther = samplesIdx["/memory/classes/metadata/other:bytes"]
	metricsMemoryClassesOsStacks = samplesIdx["/memory/classes/os-stacks:bytes"]
	metricsMemoryClassesOther = samplesIdx["/memory/classes/other:bytes"]
	metricsMemoryClassesProfilingBuckets = samplesIdx["/memory/classes/profiling/buckets:bytes"]
	metricsMemoryClassesTotal = samplesIdx["/memory/classes/total:bytes"]
	metricsSchedGoroutines = samplesIdx["/sched/goroutines:goroutines"]
	metricsSchedLatencies = samplesIdx["/sched/latencies:seconds"]
}

// maxBuckets is the maximum number of buckets we'll plots in heatmaps.
// Histograms with more buckets than that are going to be downsampled.
const maxBuckets = 100

// downsampleFactor computes the downsampling factor to use in
// downsampleHistogram, given the number of buckets in an histogram and the
// maximum number of buckets.
func downsampleFactor(nbuckets, maxbuckets int) int {
	mod := nbuckets % maxbuckets
	if mod == 0 {
		return nbuckets / maxbuckets
	}
	return 1 + nbuckets/maxbuckets
}

// downsampleBuckets downsamples the buckets in the provided histogram, using
// the given factor. The first bucket is not considered since we're only
// interested by upper bounds. If the last bucket is +Inf it gets replaced by a
// number, based on the 2 previous buckets.
func downsampleBuckets(h *metrics.Float64Histogram, factor int) []float64 {
	var ret []float64
	vals := h.Buckets[1:]

	for i := 0; i < len(vals); i++ {
		if (i+1)%factor == 0 {
			ret = append(ret, vals[i])
		}
	}
	if len(vals)%factor != 0 {
		// If the number of bucket is not divisible by the factor, let's make a
		// last downsampled bucket, even if it doesn't 'contain' the same number
		// of original buckets.
		ret = append(ret, vals[len(vals)-1])
	}

	if len(ret) > 2 && math.IsInf(ret[len(ret)-1], 1) {
		// Plotly doesn't accept an Inf bucket bound, so in this case we set the
		// last bound so that the 2 last buckets had the same size.
		ret[len(ret)-1] = ret[len(ret)-2] - ret[len(ret)-3] + ret[len(ret)-2]
	}

	return ret
}

func downsampleCounts(h *metrics.Float64Histogram, factor int) []uint64 {
	var ret []uint64
	vals := h.Counts

	if factor == 1 {
		ret = make([]uint64, len(vals))
		copy(ret, vals)
		return ret
	}

	var sum uint64
	for i := 0; i < len(vals); i++ {
		if i%factor == 0 && i > 1 {
			ret = append(ret, sum)
			sum = vals[i]
		} else {
			sum += vals[i]
		}
	}

	// Whatever sum remains, it goes to the last bucket.
	return append(ret, sum)
}

var (
	once  = sync.Once{}
	pd    *plot.Config
	am    allMetrics
	plots []plotdef
)

func plotsdef() *plot.Config {
	once.Do(createPlotsDef)
	return pd
}

func createPlotsDef() {
	am.init()

	plots = append(plots, makeHeapGlobalPlot(&am))
	plots = append(plots, makeHeapDetailsPlot(&am))
	plots = append(plots, makeLiveObjectsPlot(&am))
	plots = append(plots, makeLiveBytesPlot(&am))
	plots = append(plots, makeMSpanMCachePlot(&am))
	plots = append(plots, makeGoroutinesPlot(&am))
	plots = append(plots, makeSizeClassesPlot(&am))
	plots = append(plots, makeGCPausesPlot(&am))
	plots = append(plots, makeSchedLatPlot(&am))

	metrics.Read(am.samples)

	var layouts []interface{}
	for _, p := range plots {
		if p.isEnabled() {
			layouts = append(layouts, p.layout(am.samples))
		}
	}

	pd = &plot.Config{
		Events: []string{"lastgc"},
		Series: layouts,
	}
}

func plotsValues(samples []metrics.Sample) map[string]interface{} {
	m := make(map[string]interface{})

	for _, p := range plots {
		if p.isEnabled() {
			m[p.name()] = p.values(samples)
		}
	}

	// lastgc time series is used as source to represent garbage collection
	// timestamps as vertical bars on certain plots.
	gcStats := debug.GCStats{}
	debug.ReadGCStats(&gcStats)
	// In javascript, timestamps are in ms.
	lastgc := gcStats.LastGC.UnixMilli()
	m["lastgc"] = []int64{lastgc}
	return m
}
