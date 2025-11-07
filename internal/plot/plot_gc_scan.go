package plot

import "runtime/metrics"

var _ = register(description{
	name: "gc-scan",
	tags: []tag{tagGC},
	metrics: []string{
		"/gc/scan/globals:bytes",
		"/gc/scan/heap:bytes",
		"/gc/scan/stack:bytes",
	},
	layout: Scatter{
		Name:   "TODO(set later)",
		Title:  "GC Scan",
		Type:   "bar",
		Events: "lastgc",
		Layout: ScatterLayout{
			BarMode: "stack",
			Yaxis: ScatterYAxis{
				TickSuffix: "B",
				Title:      "bytes",
			},
		},
		Subplots: []Subplot{
			{
				Name:    "scannable globals",
				Unitfmt: "%{y:.4s}B",
				Type:    "bar",
			},
			{
				Name:    "scannable heap",
				Unitfmt: "%{y:.4s}B",
				Type:    "bar",
			},
			{
				Name:    "scanned stack",
				Unitfmt: "%{y:.4s}B",
				Type:    "bar",
			},
		},
		InfoText: `
This plot shows the amount of memory that is scannable by the GC.
<i>scannable globals</i> is <b>/gc/scan/globals</b>, the total amount of global variable space that is scannable.
<i>scannable heap</i> is <b>/gc/scan/heap</b>, the total amount of heap space that is scannable.
<i>scanned stack</i> is <b>/gc/scan/stack</b>, the number of bytes of stack that were scanned last GC cycle.
`,
	},
	make: func(idx ...int) metricsGetter {
		return &gcScan{
			idxGlobals: idx[0],
			idxHeap:    idx[1],
			idxStack:   idx[2],
		}
	},
})

type gcScan struct {
	idxGlobals int
	idxHeap    int
	idxStack   int
}

func (p *gcScan) values(samples []metrics.Sample) any {
	globals := samples[p.idxGlobals].Value.Uint64()
	heap := samples[p.idxHeap].Value.Uint64()
	stack := samples[p.idxStack].Value.Uint64()
	return []uint64{
		globals,
		heap,
		stack,
	}
}
