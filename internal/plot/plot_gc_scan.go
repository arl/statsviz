package plot

import (
	"runtime/metrics"
	"time"
)

var _ = register(description{
	metrics: []string{
		"/gc/scan/globals:bytes",
		"/gc/scan/heap:bytes",
		"/gc/scan/stack:bytes",
	},
	getvalues: func() getvalues {
		return func(_ time.Time, samples []metrics.Sample) any {
			globals := samples[idx_gc_scan_globals_bytes].Value.Uint64()
			heap := samples[idx_gc_scan_heap_bytes].Value.Uint64()
			stack := samples[idx_gc_scan_stack_bytes].Value.Uint64()

			return []uint64{
				globals,
				heap,
				stack,
			}
		}
	},
	layout: Scatter{
		Name:   "gc-scan",
		Tags:   []tag{tagGC},
		Title:  "GC Scan",
		Type:   "bar",
		Events: "lastgc",
		Layout: ScatterLayout{
			BarMode: "stack",
			Yaxis: ScatterYAxis{
				Title:      "bytes",
				TickSuffix: "B",
			},
		},
		Subplots: []Subplot{
			{Name: "scannable globals", Unitfmt: "%{y:.4s}B"},
			{Name: "scannable heap", Unitfmt: "%{y:.4s}B"},
			{Name: "scanned stack", Unitfmt: "%{y:.4s}B"},
		},
		InfoText: `<i>scannable globals</i> is <b>/gc/scan/globals:bytes</b>. The total amount of global variable space that is scannable.
<i>scannable heap</i> is <b>/gc/scan/heap:bytes</b>. The total amount of heap space that is scannable.
<i>scanned stack</i> is <b>/gc/scan/stack:bytes</b>. The number of bytes of stack that were scanned last GC cycle.`,
	},
})
