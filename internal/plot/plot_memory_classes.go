package plot

import (
	"runtime/metrics"
	"time"
)

var _ = register(description{
	metrics: []string{
		"/memory/classes/os-stacks:bytes",
		"/memory/classes/other:bytes",
		"/memory/classes/profiling/buckets:bytes",
		"/memory/classes/total:bytes",
	},
	getvalues: func() getvalues {
		return func(_ time.Time, samples []metrics.Sample) any {
			osStacks := samples[idx_memory_classes_os_stacks_bytes].Value.Uint64()
			other := samples[idx_memory_classes_other_bytes].Value.Uint64()
			profBuckets := samples[idx_memory_classes_profiling_buckets_bytes].Value.Uint64()
			total := samples[idx_memory_classes_total_bytes].Value.Uint64()

			return []uint64{
				osStacks,
				other,
				profBuckets,
				total,
			}
		}
	},
	layout: Scatter{
		Name:   "memory-classes",
		Tags:   []tag{tagGC},
		Title:  "Memory classes",
		Type:   "scatter",
		Events: "lastgc",
		Layout: ScatterLayout{
			Yaxis: ScatterYAxis{
				Title:      "bytes",
				TickSuffix: "B",
			},
		},
		Subplots: []Subplot{
			{Name: "os stacks", Unitfmt: "%{y:.4s}B"},
			{Name: "other", Unitfmt: "%{y:.4s}B"},
			{Name: "profiling buckets", Unitfmt: "%{y:.4s}B"},
			{Name: "total", Unitfmt: "%{y:.4s}B"},
		},
		InfoText: `<i>os stacks</i> is <b>/memory/classes/os-stacks:bytes</b>. Stack memory allocated by the underlying operating system.
<i>other</i> is <b>/memory/classes/other:bytes</b>. Memory used by execution trace buffers, structures for debugging the runtime, finalizer and profiler specials, and more.
<i>profiling buckets</i> is <b>/memory/classes/profiling/buckets:bytes</b>. Memory that is used by the stack trace hash map used for profiling.
<i>total</i> is <b>/memory/classes/total:bytes</b>. All memory mapped by the Go runtime into the current process as read-write.`,
	},
})
