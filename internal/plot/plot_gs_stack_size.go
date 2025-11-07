package plot

import "runtime/metrics"

type gcStackSize struct {
	idxstack int
}

func makeGCStackSize(indices ...int) metricsGetter {
	return &gcStackSize{
		idxstack: indices[0],
	}
}

func (p *gcStackSize) values(samples []metrics.Sample) any {
	stackSize := samples[p.idxstack].Value.Uint64()
	return []uint64{stackSize}
}

var gcStackSizeLayout = Scatter{
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
}
