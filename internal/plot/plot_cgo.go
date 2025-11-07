package plot

import (
	"runtime/metrics"
	"time"
)

var _ = register(description{
	name: "cgo",
	tags: []tag{tagMisc},
	metrics: []string{
		"/cgo/go-to-c-calls:calls",
	},
	getvalues: func() getvalues {
		// TODO show cgo calls per second
		deltago2c := deltaUint64(idxcgogotocalls)

		return func(_ time.Time, samples []metrics.Sample) any {
			return []uint64{deltago2c(samples)}
		}
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
})
