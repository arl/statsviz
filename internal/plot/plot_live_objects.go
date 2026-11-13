package plot

import (
	"runtime/metrics"
	"time"
)

var _ = register(description{
	metrics: []string{
		"/gc/heap/objects:objects",
	},
	getvalues: func() getvalues {
		return func(_ time.Time, samples []metrics.Sample) any {
			gcHeapObjects := samples[idx_gc_heap_objects_objects].Value.Uint64()

			return []uint64{
				gcHeapObjects,
			}
		}
	},
	layout: Scatter{
		Name:   "live-objects",
		Tags:   []tag{tagGC},
		Title:  "Live Objects in Heap",
		Type:   "bar",
		Events: "lastgc",
		Layout: ScatterLayout{
			BarMode: "stack",
			Yaxis: ScatterYAxis{
				Title:   "objects",
				Integer: true,
			},
		},
		Subplots: []Subplot{
			{Name: "live objects", Unitfmt: "%{y}"},
		},
		InfoText: `<i>live objects</i> is <b>/gc/heap/objects:objects</b>. Number of objects, live or unswept, occupying heap memory.`,
	},
})
