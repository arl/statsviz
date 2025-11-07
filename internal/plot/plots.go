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

type plotDesc struct {
	name    string
	tags    []tag
	metrics []string
	layout  any

	// make creates the state (support struct) for the plot.
	make func(indices ...int) metricsGetter
}

var (
	plotDescs []plotDesc

	metricDescs = metrics.All()
	metricIdx   map[string]int
)

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

	plotDescs = []plotDesc{
		{
			name: "garbage collection",
			tags: []tag{tagGC},
			metrics: []string{
				"/gc/gomemlimit:bytes",
				"/gc/heap/live:bytes",
				"/gc/heap/goal:bytes",
				"/memory/classes/total:bytes",
				"/memory/classes/heap/released:bytes",
			},
			layout: garbageCollectionLayout,
			make:   makeGarbageCollection,
		},
		{
			name: "heap (details)",
			tags: []tag{tagGC},
			metrics: []string{
				"/memory/classes/heap/objects:bytes",
				"/memory/classes/heap/unused:bytes",
				"/memory/classes/heap/free:bytes",
				"/memory/classes/heap/released:bytes",
				"/memory/classes/heap/stacks:bytes",
				"/gc/heap/goal:bytes",
			},
			layout: heapDetailslLayout,
			make:   makeHeapDetails,
		},
		{
			name: "live-objects",
			tags: []tag{tagGC},
			metrics: []string{
				"/gc/heap/objects:objects",
			},
			layout: liveObjectsLayout,
			make:   makeLiveObjects,
		},
		{
			name: "live-bytes",
			tags: []tag{tagGC},
			metrics: []string{
				"/gc/heap/allocs:bytes",
				"/gc/heap/frees:bytes",
			},
			layout: liveBytesLayout,
			make:   makeLiveBytes,
		},
		{
			name: "mspan-mcache",
			tags: []tag{tagGC},
			metrics: []string{
				"/memory/classes/metadata/mspan/inuse:bytes",
				"/memory/classes/metadata/mspan/free:bytes",
				"/memory/classes/metadata/mcache/inuse:bytes",
				"/memory/classes/metadata/mcache/free:bytes",
			},
			layout: mspanMCacheLayout,
			make:   makeMSpanMCache,
		},
		{
			name: "size-classes",
			tags: []tag{tagGC},
			metrics: []string{
				"/gc/heap/allocs-by-size:bytes",
				"/gc/heap/frees-by-size:bytes",
			},
			layout: sizeClassesLayout(samples),
			make:   makeSizeClasses,
		},
		{
			name: "runnable-time",
			tags: []tag{tagScheduler},
			metrics: []string{
				"/sched/latencies:seconds",
			},
			layout: runnableTimeLayout(samples),
			make:   makeRunnableTime,
		},
		{
			name: "sched-events",
			tags: []tag{tagScheduler},
			metrics: []string{
				"/sched/latencies:seconds",
				"/sched/gomaxprocs:threads",
			},
			layout: schedEventsLayout,
			make:   makeSchedEvents,
		},
		{
			name: "cgo",
			tags: []tag{tagMisc},
			metrics: []string{
				"/cgo/go-to-c-calls:calls",
			},
			layout: cgoLayout,
			make:   makeCGO,
		},
		{
			name: "gc-stack-size",
			tags: []tag{tagGC},
			metrics: []string{
				"/gc/stack/starting-size:bytes",
			},
			layout: gcStackSizeLayout,
			make:   makeGCStackSize,
		},
		{
			name: "gc-cycles",
			tags: []tag{tagGC},
			metrics: []string{
				"/gc/cycles/automatic:gc-cycles",
				"/gc/cycles/forced:gc-cycles",
				"/gc/cycles/total:gc-cycles",
			},
			layout: gcCyclesLayout,
			make:   makeGCCycles,
		},
		{
			name: "memory-classes",
			tags: []tag{tagGC},
			metrics: []string{
				"/memory/classes/os-stacks:bytes",
				"/memory/classes/other:bytes",
				"/memory/classes/profiling/buckets:bytes",
				"/memory/classes/total:bytes",
			},
			layout: memoryClassesLayout,
			make:   makeMemoryClasses,
		},
		{
			name: "cpu-gc",
			tags: []tag{tagCPU, tagGC},
			metrics: []string{
				"/cpu/classes/gc/mark/assist:cpu-seconds",
				"/cpu/classes/gc/mark/dedicated:cpu-seconds",
				"/cpu/classes/gc/mark/idle:cpu-seconds",
				"/cpu/classes/gc/pause:cpu-seconds",
			},
			layout: cpuGCLayout,
			make:   makeCPUgc,
		},
		{
			name: "cpu-scavenger",
			tags: []tag{tagCPU, tagGC},
			metrics: []string{
				"/cpu/classes/scavenge/assist:cpu-seconds",
				"/cpu/classes/scavenge/background:cpu-seconds",
			},
			layout: cpuScavengerLayout,
			make:   makeCPUscavenger,
		},
		{
			name: "cpu-overall",
			tags: []tag{tagCPU},
			metrics: []string{
				"/cpu/classes/user:cpu-seconds",
				"/cpu/classes/scavenge/total:cpu-seconds",
				"/cpu/classes/idle:cpu-seconds",
				"/cpu/classes/gc/total:cpu-seconds",
				"/cpu/classes/total:cpu-seconds",
			},
			layout: cpuOverallLayout,
			make:   makeCPUoverall,
		},
		{
			name: "mutex-wait",
			tags: []tag{tagMisc},
			metrics: []string{
				"/sync/mutex/wait/total:seconds",
			},
			layout: mutexWaitLayout,
			make:   makeMutexWait,
		},
		{
			name: "gc-scan",
			tags: []tag{tagGC},
			metrics: []string{
				"/gc/scan/globals:bytes",
				"/gc/scan/heap:bytes",
				"/gc/scan/stack:bytes",
			},
			layout: gcScanLayout,
			make:   makeGCScan,
		},
		{
			name: "alloc-free-rate",
			tags: []tag{tagGC},
			metrics: []string{
				"/gc/heap/allocs:objects",
				"/gc/heap/frees:objects",
			},
			layout: allocFreeRatesLayout,
			make:   makeAllocFreeRates,
		},
		{
			name: "total-pauses-gc",
			tags: []tag{tagScheduler},
			metrics: []string{
				"/sched/pauses/total/gc:seconds",
			},
			layout: totalPausesGCLayout(samples),
			make:   makeTotalPausesGC,
		},
		{
			name: "total-pauses-other",
			tags: []tag{tagScheduler},
			metrics: []string{
				"/sched/pauses/total/other:seconds",
			},
			layout: totalPausesOtherLayout(samples),
			make:   makeTotalPausesOther,
		},
		{
			name: "stopping-pauses-gc",
			tags: []tag{tagScheduler},
			metrics: []string{
				"/sched/pauses/stopping/gc:seconds",
			},
			layout: stoppingPausesGCLayout(samples),
			make:   makeStoppingPausesGC,
		},
		{
			name: "stopping-pauses-other",
			tags: []tag{tagScheduler},
			metrics: []string{
				"/sched/pauses/stopping/other:seconds",
			},
			layout: stoppingPausesOtherLayout(samples),
			make:   makeStoppingPausesOther,
		},
		{
			name: "goroutines",
			tags: []tag{tagScheduler},
			metrics: []string{
				"/sched/goroutines:goroutines",
				"/sched/goroutines-created:goroutines",
				"/sched/goroutines/not-in-go:goroutines",
				"/sched/goroutines/runnable:goroutines",
				"/sched/goroutines/running:goroutines",
				"/sched/goroutines/waiting:goroutines",
			},
			layout: goroutinesLayout,
			make:   makeGoroutines,
		},
		{
			name: "threads",
			tags: []tag{tagScheduler},
			metrics: []string{
				"/sched/threads/total:threads",
			},
			layout: threadsLayout,
			make:   makeThreads,
		},
	}
}
