package plot

import "runtime/metrics"

type liveBytes struct {
	idxallocs int
	idxfrees  int
}

func makeLiveBytes(indices ...int) metricsGetter {
	return &liveBytes{
		idxallocs: indices[0],
		idxfrees:  indices[1],
	}
}

func (p *liveBytes) values(samples []metrics.Sample) any {
	allocBytes := samples[p.idxallocs].Value.Uint64()
	freedBytes := samples[p.idxfrees].Value.Uint64()
	return []uint64{
		allocBytes - freedBytes,
	}
}

var liveBytesLayout = Scatter{
	Name:   "TODO(set later)",
	Title:  "Live Bytes in Heap",
	Type:   "bar",
	Events: "lastgc",
	Layout: ScatterLayout{
		Yaxis: ScatterYAxis{
			Title: "bytes",
		},
	},
	Subplots: []Subplot{
		{
			Name:    "live bytes",
			Unitfmt: "%{y:.4s}B",
			Color:   RGBString(135, 182, 218),
		},
	},
	InfoText: `<i>Live bytes</i> is <b>/gc/heap/allocs - /gc/heap/frees</b>. It's the number of bytes currently allocated (and not yet GC'ec) to the heap by the application.`,
}
