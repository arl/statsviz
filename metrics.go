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
	once = sync.Once{}
	pd   *plot.Definition

	gcpausesFactor int
	schedlatFactor int
)

func plotsdef() *plot.Definition {
	once.Do(createPlotsDef)
	return pd
}

func createPlotsDef() {
	// Sample the metric once
	metrics.Read(samples)

	// TODO(arl) rename metrics so that they match that of the new package (example: nextGC -> Gc heap goal)

	// Perform a sanity check on the number of buckets on the 'allocs' and
	// 'frees' size classes histograms. Statsviz plots a single histogram based
	// on those 2 so we want them to have the same number of buckets, which
	// should be true.
	allocsBySize := samples[metricGCHeapAllocsBySize].Value.Float64Histogram()
	freesBySize := samples[metricGCHeapFreesBySize].Value.Float64Histogram()
	if len(allocsBySize.Buckets) != len(freesBySize.Buckets) {
		panic("different number of buckets in allocs and frees size classes histograms!")
	}

	// No downsampling for the size classes histogram (factor=1) but we still
	// need to adapt boundaries to plotly heatmaps.
	sizeClassesBuckets := downsampleBuckets(allocsBySize, 1)

	gcpauses := samples[metricsGCPauses].Value.Float64Histogram()
	gcpausesFactor = downsampleFactor(len(gcpauses.Buckets), maxBuckets)
	gcpausesBuckets := downsampleBuckets(gcpauses, gcpausesFactor)

	schedlat := samples[metricsSchedLatencies].Value.Float64Histogram()
	schedlatFactor = downsampleFactor(len(schedlat.Buckets), maxBuckets)
	schedlatBuckets := downsampleBuckets(schedlat, schedlatFactor)

	pd = &plot.Definition{
		Events: []string{"lastgc"},
		Series: []interface{}{
			plot.Scatter{
				Name:       "heap-global",
				Title:      "Heap (global)",
				Type:       "scatter",
				HorzEvents: "lastgc",
				Layout: plot.ScatterLayout{
					Yaxis: plot.ScatterLayoutYAxis{
						Title:      "bytes",
						TickSuffix: "B",
					},
				},
				Subplots: []plot.Subplot{
					{
						Name:       "heap in-use",
						Unitfmt:    "%{y:.4s}B",
						HoverOn:    "points+fills",
						StackGroup: "one",
					},
					{
						Name:       "heap free",
						Unitfmt:    "%{y:.4s}B",
						HoverOn:    "points+fills",
						StackGroup: "one",
					},
					{
						Name:       "heap released",
						Unitfmt:    "%{y:.4s}B",
						HoverOn:    "points+fills",
						StackGroup: "one",
					},
				},
			},

			plot.Scatter{
				Name:       "heap-details",
				Title:      "Heap (details)",
				Type:       "scatter",
				HorzEvents: "lastgc",
				Layout: plot.ScatterLayout{
					Yaxis: plot.ScatterLayoutYAxis{
						Title:      "bytes",
						TickSuffix: "B",
					},
				},
				Subplots: []plot.Subplot{
					{
						Name:    "heap sys",
						Unitfmt: "%{y:.4s}B",
					},
					{
						Name:    "heap objects",
						Unitfmt: "%{y:.4s}B",
					},
					{
						Name:    "heap stacks",
						Unitfmt: "%{y:.4s}B",
					},
					{
						Name:    "heap goal",
						Unitfmt: "%{y:.4s}B",
					},
				},
			},

			plot.Scatter{
				Name:       "live bytes",
				Title:      "Live Bytes in Heap",
				Type:       "bar",
				HorzEvents: "lastgc",
				Layout: plot.ScatterLayout{
					Yaxis: plot.ScatterLayoutYAxis{
						Title: "bytes",
					},
				},
				Subplots: []plot.Subplot{
					{
						Name:    "live bytes",
						Unitfmt: "%{y:.4s}B",
						Color:   plot.RGBString(135, 182, 218),
					},
				},
			},

			plot.Scatter{
				Name:       "live objects",
				Title:      "Live Objects in Heap",
				Type:       "bar",
				HorzEvents: "lastgc",
				Layout: plot.ScatterLayout{
					Yaxis: plot.ScatterLayoutYAxis{
						Title: "objects",
					},
				},
				Subplots: []plot.Subplot{
					{
						Name:    "live objects",
						Unitfmt: "%{y:.4s}",
						Color:   plot.RGBString(255, 195, 128),
					},
				},
			},
			plot.Scatter{
				Name:       "mspan-mcache",
				Title:      "MSpan/MCache",
				Type:       "scatter",
				HorzEvents: "lastgc",
				Layout: plot.ScatterLayout{
					Yaxis: plot.ScatterLayoutYAxis{
						Title:      "bytes",
						TickSuffix: "B",
					},
				},
				Subplots: []plot.Subplot{
					{
						Name:    "mspan in-use",
						Unitfmt: "%{y:.4s}B",
					},
					{
						Name:    "mspan sys",
						Unitfmt: "%{y:.4s}B",
					},
					{
						Name:    "mcache in-use",
						Unitfmt: "%{y:.4s}B",
					},
					{
						Name:    "mcache sys",
						Unitfmt: "%{y:.4s}B",
					},
				},
			},
			plot.Scatter{
				Name:       "goroutines",
				Title:      "Goroutines",
				Type:       "scatter",
				HorzEvents: "lastgc",
				Layout: plot.ScatterLayout{
					Yaxis: plot.ScatterLayoutYAxis{
						Title: "goroutines",
					},
				},
				Subplots: []plot.Subplot{
					{
						Name:    "goroutines",
						Unitfmt: "%{y}",
					},
				},
			},
			plot.Heatmap{
				Name:       "sizeclasses",
				Title:      "Size Classes",
				Type:       "heatmap",
				UpdateFreq: 5,
				HorzEvents: "",
				Layout: plot.HeatmapLayout{
					Yaxis: plot.HeatmapLayoutYAxis{
						Title: "size class",
					},
				},
				Colorscale: plot.BlueShades,
				Buckets:    floatseq(len(sizeClassesBuckets)),
				CustomData: sizeClassesBuckets,
				Hover: plot.HeapmapHover{
					YName: "size class",
					YUnit: "bytes",
					ZName: "objects",
				},
			},
			plot.Heatmap{
				Name:       "gcpauses",
				Title:      "Stop-the-world pause latencies",
				Type:       "heatmap",
				UpdateFreq: 5,
				HorzEvents: "",
				Layout: plot.HeatmapLayout{
					Yaxis: plot.HeatmapLayoutYAxis{
						Title: "pause duration",
					},
				},
				Colorscale: plot.PinkShades,
				Buckets:    floatseq(len(gcpausesBuckets)),
				CustomData: gcpausesBuckets,
				Hover: plot.HeapmapHover{
					YName: "pause duration",
					YUnit: "duration",
					ZName: "pauses",
				},
			},
			plot.Heatmap{
				Name:       "sched-latencies",
				Title:      "Time in scheduler before a goroutine runs",
				Type:       "heatmap",
				UpdateFreq: 5,
				HorzEvents: "",
				Layout: plot.HeatmapLayout{
					Yaxis: plot.HeatmapLayoutYAxis{
						Title: "duration",
					},
				},
				Colorscale: plot.GreenShades,
				Buckets:    floatseq(len(schedlatBuckets)),
				CustomData: schedlatBuckets,
				Hover: plot.HeapmapHover{
					YName: "duration",
					YUnit: "duration",
					ZName: "goroutines",
				},
			},
		},
	}
}

func plotsValues(samples []metrics.Sample) map[string]interface{} {
	m := make(map[string]interface{})

	heapObjects := samples[metricsMemoryClassesHeapObjects].Value.Uint64()
	heapUnused := samples[metricsMemoryClassesHeapUnusedBytes].Value.Uint64()
	heapInUse := heapObjects + heapUnused

	heapFree := samples[metricsMemoryClassesHeapFreeBytes].Value.Uint64()
	heapReleased := samples[metricsMemoryClassesHeapReleasedBytes].Value.Uint64()

	m["heap-global"] = []uint64{
		heapInUse,
		heapFree,
		heapReleased,
	}

	heapIdle := heapReleased + heapFree
	heapSys := heapInUse + heapIdle
	heapStacks := samples[metricsMemoryClassesHeapStackBytes].Value.Uint64()
	nextGC := samples[metricsGCHeapGoalBytes].Value.Uint64()

	m["heap-details"] = []uint64{
		heapSys,
		heapObjects,
		heapStacks,
		nextGC,
	}

	gcHeapObjects := samples[metricsGCHeapObjects].Value.Uint64()
	m["live objects"] = []uint64{
		gcHeapObjects,
	}

	allocBytes := samples[metricGCHeapAllocsBytes].Value.Uint64()
	freedBytes := samples[metricGCHeapFreesBytes].Value.Uint64()
	m["live bytes"] = []uint64{
		allocBytes - freedBytes,
	}

	mspanInUse := samples[metricsMemoryClassesMetadataMSpanInUse].Value.Uint64()
	mspanSys := samples[metricsMemoryClassesMetadataMSpanFree].Value.Uint64()
	mcacheInUse := samples[metricsMemoryClassesMetadataMCacheInUse].Value.Uint64()
	mcacheSys := samples[metricsMemoryClassesMetadataMCacheFree].Value.Uint64()
	m["mspan-mcache"] = []uint64{
		mspanInUse,
		mspanSys,
		mcacheInUse,
		mcacheSys,
	}

	m["goroutines"] = []uint64{samples[metricsSchedGoroutines].Value.Uint64()}

	// Now we take lastGC from GCstats
	gcStats := debug.GCStats{}
	debug.ReadGCStats(&gcStats)
	// Javascript datetime is in ms
	m["lastgc"] = []int64{gcStats.LastGC.UnixMilli()}

	allocsBySize := samples[metricGCHeapAllocsBySize].Value.Float64Histogram()
	freesBySize := samples[metricGCHeapFreesBySize].Value.Float64Histogram()
	sizeClasses := make([]uint64, len(allocsBySize.Counts))
	for i := 0; i < len(sizeClasses); i++ {
		sizeClasses[i] = allocsBySize.Counts[i] - freesBySize.Counts[i]
	}
	m["sizeclasses"] = sizeClasses

	gcpauses := samples[metricsGCPauses].Value.Float64Histogram()
	m["gcpauses"] = downsampleCounts(gcpauses, gcpausesFactor)

	schedlat := samples[metricsSchedLatencies].Value.Float64Histogram()
	m["sched-latencies"] = downsampleCounts(schedlat, schedlatFactor)
	return m
}

func floatseq(n int) []float64 {
	seq := make([]float64, n)
	for i := 0; i < n; i++ {
		seq[i] = float64(i)
	}
	return seq
}
