package plot

import "runtime/metrics"

var _ = register(description{
	name: "live-bytes",
	tags: []tag{tagGC},
	metrics: []string{
		"/gc/heap/allocs:bytes",
		"/gc/heap/frees:bytes",
	},
	layout: Scatter{
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
	},
	make: func(idx ...int) metricsGetter {
		return &liveBytes{
			idxallocs: idx[0],
			idxfrees:  idx[1],
		}
	},
})

type liveBytes struct {
	idxallocs int
	idxfrees  int
}

func (p *liveBytes) values(samples []metrics.Sample) any {
	allocBytes := samples[p.idxallocs].Value.Uint64()
	freedBytes := samples[p.idxfrees].Value.Uint64()
	return []uint64{
		allocBytes - freedBytes,
	}
}
