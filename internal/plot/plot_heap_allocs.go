package plot

import (
	"runtime/metrics"
	"time"
)

var _ = register(description{
	metrics: []string{
		"/gc/heap/allocs:objects",
		"/gc/heap/frees:objects",
		"/gc/cycles/total:gc-cycles",
	},
	getvalues: func() getvalues {
		rateallocs := gcRate[uint64]()
		ratefrees := gcRate[uint64]()

		return func(now time.Time, samples []metrics.Sample) any {
			allocs := samples[idx_gc_heap_allocs_objects].Value.Uint64()
			frees := samples[idx_gc_heap_frees_objects].Value.Uint64()
			cycles := samples[idx_gc_cycles_total_gc_cycles].Value.Uint64()

			return []float64{
				rateallocs(now, cycles, allocs),
				ratefrees(now, cycles, frees),
			}
		}
	},
	layout: Scatter{
		Name:   "alloc-free-rate",
		Tags:   []tag{tagGC},
		Title:  "Heap Allocation & Free Rates",
		Type:   "scatter",
		Events: "lastgc",
		Layout: ScatterLayout{
			Yaxis: ScatterYAxis{
				Title: "objects / second",
			},
		},
		Subplots: []Subplot{
			{Name: "allocs/sec", Unitfmt: "%{y:.4s}"},
			{Name: "frees/sec", Unitfmt: "%{y:.4s}"},
		},
		InfoText: `<i>allocs/sec</i> is the rate of change, per second, averaged over the last GC cycle, of <b>/gc/heap/allocs:objects</b>. Cumulative count of heap allocations triggered by the application.
<i>frees/sec</i> is the rate of change, per second, averaged over the last GC cycle, of <b>/gc/heap/frees:objects</b>. Cumulative count of heap allocations whose storage was freed by the garbage collector.`,
	},
})
