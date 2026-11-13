package plot

import (
	"runtime/metrics"
	"time"
)

var _ = register(description{
	metrics: []string{
		"/gc/heap/allocs:bytes",
		"/gc/heap/frees:bytes",
	},
	getvalues: func() getvalues {
		return func(_ time.Time, samples []metrics.Sample) any {
			allocBytes := samples[idx_gc_heap_allocs_bytes].Value.Uint64()
			freedBytes := samples[idx_gc_heap_frees_bytes].Value.Uint64()

			return []uint64{
				allocBytes - freedBytes,
			}
		}
	},
	layout: Scatter{
		Name:   "live-bytes",
		Tags:   []tag{tagGC},
		Title:  "Live Bytes in Heap",
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
			{Name: "live bytes", Unitfmt: "%{y:.4s}B"},
		},
		InfoText: `<i>live bytes</i> is <b>/gc/heap/allocs:bytes - /gc/heap/frees:bytes</b>.`,
	},
})
