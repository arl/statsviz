package plot

import "runtime/metrics"

type gcScan struct {
	idxGlobals int
	idxHeap    int
	idxStack   int
}

func makeGCScan(indices ...int) metricsGetter {
	return &gcScan{
		idxGlobals: indices[0],
		idxHeap:    indices[1],
		idxStack:   indices[2],
	}
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

var gcScanLayout = Scatter{
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
}
