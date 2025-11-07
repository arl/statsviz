package plot

import (
	"runtime/metrics"
	"time"
)

var _ = register(description{
	name: "alloc-free-rate",
	tags: []tag{tagGC},
	metrics: []string{
		"/gc/heap/allocs:objects",
		"/gc/heap/frees:objects",
	},
	layout: Scatter{
		Name:  "heap alloc/free rates",
		Title: "Heap Allocation & Free Rates",
		Type:  "scatter",
		Layout: ScatterLayout{
			Yaxis: ScatterYAxis{
				Title: "objects / second",
			},
		},
		Subplots: []Subplot{
			{
				Name:    "allocs/sec",
				Unitfmt: "%{y:.4s}",
				Color:   RGBString(66, 133, 244),
			},
			{
				Name:    "frees/sec",
				Unitfmt: "%{y:.4s}",
				Color:   RGBString(219, 68, 55),
			},
		},
		InfoText: `
<i>Allocations per second</i> is derived by differencing the cumulative <b>/gc/heap/allocs:objects</b> metric.
<i>Frees per second</i> is similarly derived from <b>/gc/heap/frees:objects</b>.`,
	},
	make: func(idx ...int) metricsGetter {
		return &allocFreeRates{
			idxallocs: idx[0],
			idxfrees:  idx[1],
		}
	},
})

type allocFreeRates struct {
	idxallocs int
	idxfrees  int

	lasttime   time.Time
	lastallocs uint64
	lastfrees  uint64
}

func (p *allocFreeRates) values(samples []metrics.Sample) any {
	if p.lasttime.IsZero() {
		p.lasttime = time.Now()
		p.lastallocs = samples[p.idxallocs].Value.Uint64()
		p.lastfrees = samples[p.idxfrees].Value.Uint64()

		return []float64{0, 0}
	}

	t := time.Since(p.lasttime).Seconds()

	allocs := float64(samples[p.idxallocs].Value.Uint64()-p.lastallocs) / t
	frees := float64(samples[p.idxfrees].Value.Uint64()-p.lastfrees) / t

	p.lastallocs = samples[p.idxallocs].Value.Uint64()
	p.lastfrees = samples[p.idxfrees].Value.Uint64()
	p.lasttime = time.Now()

	return []float64{
		allocs,
		frees,
	}
}
