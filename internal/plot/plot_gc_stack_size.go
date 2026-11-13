package plot

import (
	"runtime/metrics"
	"time"
)

var _ = register(description{
	metrics: []string{
		"/gc/stack/starting-size:bytes",
	},
	getvalues: func() getvalues {
		return func(_ time.Time, samples []metrics.Sample) any {
			stackSize := samples[idx_gc_stack_starting_size_bytes].Value.Uint64()

			return []uint64{
				stackSize,
			}
		}
	},
	layout: Scatter{
		Name:   "gc-stack-size",
		Tags:   []tag{tagGC},
		Title:  "Goroutines stack starting size",
		Type:   "scatter",
		Events: "lastgc",
		Layout: ScatterLayout{
			Yaxis: ScatterYAxis{
				Title:      "bytes",
				TickSuffix: "B",
			},
		},
		Subplots: []Subplot{
			{Name: "new goroutines stack size", Unitfmt: "%{y:.4s}B"},
		},
		InfoText: `<i>new goroutines stack size</i> is <b>/gc/stack/starting-size:bytes</b>. The stack size of new goroutines.`,
	},
})
