package plot

import (
	"math"
	"runtime/metrics"
)

var _ = register(description{
	name: "cgo",
	tags: []tag{tagMisc},
	metrics: []string{
		"/cgo/go-to-c-calls:calls",
	},
	layout: Scatter{
		Name:  "TODO(set later)",
		Title: "CGO Calls",
		Type:  "bar",
		Layout: ScatterLayout{
			Yaxis: ScatterYAxis{
				Title: "calls",
			},
		},
		Subplots: []Subplot{
			{
				Name:    "calls from go to c",
				Unitfmt: "%{y}",
				Color:   "red",
			},
		},
		InfoText: "Shows the count of calls made from Go to C by the current process, per unit of time. Uses <b>/cgo/go-to-c-calls:calls</b>",
	},
	make: func(indices ...int) metricsGetter {
		return &cgo{
			idxgo2c:  indices[0],
			lastgo2c: math.MaxUint64,
		}
	},
})

type cgo struct {
	idxgo2c  int
	lastgo2c uint64
}

// TODO show cgo calls per second
func (p *cgo) values(samples []metrics.Sample) any {
	go2c := samples[p.idxgo2c].Value.Uint64()
	curgo2c := go2c - p.lastgo2c
	if p.lastgo2c == math.MaxUint64 {
		curgo2c = 0
	}
	p.lastgo2c = go2c

	return []uint64{curgo2c}
}
