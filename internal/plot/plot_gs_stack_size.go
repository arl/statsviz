package plot

import "runtime/metrics"

var _ = register(description{
	name: "gc-stack-size",
	tags: []tag{tagGC},
	metrics: []string{
		"/gc/stack/starting-size:bytes",
	},
	layout: Scatter{
		Name:  "TODO(set later)",
		Title: "Goroutines stack starting size",
		Type:  "scatter",
		Layout: ScatterLayout{
			Yaxis: ScatterYAxis{
				Title: "bytes",
			},
		},
		Subplots: []Subplot{
			{
				Name:    "new goroutines stack size",
				Unitfmt: "%{y:.4s}B",
			},
		},
		InfoText: "Shows the stack size of new goroutines, uses <b>/gc/stack/starting-size:bytes</b>",
	},
	make: func(indices ...int) metricsGetter {
		return &gcStackSize{
			idxstack: indices[0],
		}
	},
})

type gcStackSize struct {
	idxstack int
}

func (p *gcStackSize) values(samples []metrics.Sample) any {
	stackSize := samples[p.idxstack].Value.Uint64()
	return []uint64{stackSize}
}
