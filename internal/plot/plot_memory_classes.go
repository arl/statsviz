package plot

import "runtime/metrics"

var _ = register(description{
	name: "memory-classes",
	tags: []tag{tagGC},
	metrics: []string{
		"/memory/classes/os-stacks:bytes",
		"/memory/classes/other:bytes",
		"/memory/classes/profiling/buckets:bytes",
		"/memory/classes/total:bytes",
	},
	layout: Scatter{
		Name:   "TODO(set later)",
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
			{
				Name:    "os stacks",
				Unitfmt: "%{y:.4s}B",
			},
			{
				Name:    "other",
				Unitfmt: "%{y:.4s}B",
			},
			{
				Name:    "profiling buckets",
				Unitfmt: "%{y:.4s}B",
			},
			{
				Name:    "total",
				Unitfmt: "%{y:.4s}B",
			},
		},

		InfoText: `
<i>OS stacks</i> is <b>/memory/classes/os-stacks</b>, stack memory allocated by the underlying operating system.
<i>Other</i> is <b>/memory/classes/other</b>, memory used by execution trace buffers, structures for debugging the runtime, finalizer and profiler specials, and more.
<i>Profiling buckets</i> is <b>/memory/classes/profiling/buckets</b>, memory that is used by the stack trace hash map used for profiling.
<i>Total</i> is <b>/memory/classes/total</b>, all memory mapped by the Go runtime into the current process as read-write.`,
	},
	make: func(idx ...int) metricsGetter {
		return &memoryClasses{
			idxOSStacks:    idx[0],
			idxOther:       idx[1],
			idxProfBuckets: idx[2],
			idxTotal:       idx[3],
		}
	},
})

type memoryClasses struct {
	idxOSStacks    int
	idxOther       int
	idxProfBuckets int
	idxTotal       int
}

func (p *memoryClasses) values(samples []metrics.Sample) any {
	osStacks := samples[p.idxOSStacks].Value.Uint64()
	other := samples[p.idxOther].Value.Uint64()
	profBuckets := samples[p.idxProfBuckets].Value.Uint64()
	total := samples[p.idxTotal].Value.Uint64()

	return []uint64{
		osStacks,
		other,
		profBuckets,
		total,
	}
}
