package plot

import "runtime/metrics"

var _ = register(description{
	name: "live-objects",
	tags: []tag{tagGC},
	metrics: []string{
		"/gc/heap/objects:objects",
	},
	layout: Scatter{
		Name:   "TODO(set later)",
		Title:  "Live Objects in Heap",
		Type:   "bar",
		Events: "lastgc",
		Layout: ScatterLayout{
			Yaxis: ScatterYAxis{
				Title: "objects",
			},
		},
		Subplots: []Subplot{
			{
				Name:    "live objects",
				Unitfmt: "%{y:.4s}",
				Color:   RGBString(255, 195, 128),
			},
		},
		InfoText: `<i>Live objects</i> is <b>/gc/heap/objects</b>. It's the number of objects, live or unswept, occupying heap memory.`,
	},
	make: func(indices ...int) metricsGetter {
		return &liveObjects{
			idxobjects: indices[0],
		}
	},
})

type liveObjects struct {
	idxobjects int
}

func (p *liveObjects) values(samples []metrics.Sample) any {
	gcHeapObjects := samples[p.idxobjects].Value.Uint64()
	return []uint64{
		gcHeapObjects,
	}
}
