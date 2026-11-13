package plot

import (
	"math"
	"runtime/metrics"
	"time"
)

var _ = register(description{
	metrics: []string{
		"/gc/gomemlimit:bytes",
		"/gc/heap/live:bytes",
		"/gc/heap/goal:bytes",
		"/memory/classes/total:bytes",
		"/memory/classes/heap/released:bytes",
	},
	getvalues: func() getvalues {
		return func(_ time.Time, samples []metrics.Sample) any {
			memLimit := samples[idx_gc_gomemlimit_bytes].Value.Uint64()
			heapLive := samples[idx_gc_heap_live_bytes].Value.Uint64()
			heapGoal := samples[idx_gc_heap_goal_bytes].Value.Uint64()
			memTotal := samples[idx_memory_classes_total_bytes].Value.Uint64()
			heapReleased := samples[idx_memory_classes_heap_released_bytes].Value.Uint64()

			if memLimit == math.MaxInt64 {
				memLimit = 0
			}

			return []uint64{
				memLimit,
				memTotal - heapReleased,
				heapLive,
				heapGoal,
			}
		}
	},
	layout: Scatter{
		Name:   "garbage collection",
		Tags:   []tag{tagGC},
		Title:  "GC Memory Summary",
		Type:   "scatter",
		Events: "lastgc",
		Layout: ScatterLayout{
			Yaxis: ScatterYAxis{
				Title:      "bytes",
				TickSuffix: "B",
			},
		},
		Subplots: []Subplot{
			{Name: "memory limit", Unitfmt: "%{y:.4s}B"},
			{Name: "in-use memory", Unitfmt: "%{y:.4s}B"},
			{Name: "heap live", Unitfmt: "%{y:.4s}B"},
			{Name: "heap goal", Unitfmt: "%{y:.4s}B"},
		},
		InfoText: `<i>memory limit</i> is <b>/gc/gomemlimit:bytes</b>. Go runtime memory limit configured by the user, otherwise 0.
<i>in-use memory</i> is <b>/memory/classes/total:bytes - /memory/classes/heap/released:bytes</b>.
<i>heap live</i> is <b>/gc/heap/live:bytes</b>. Heap memory occupied by live objects that were marked by the previous GC.
<i>heap goal</i> is <b>/gc/heap/goal:bytes</b>. Heap size target for the end of the GC cycle.`,
	},
})
