package plot

import (
	"runtime/metrics"
	"time"
)

var _ = register(description{
	metrics: []string{
		"/cgo/go-to-c-calls:calls",
	},
	getvalues: func() getvalues {
		// TODO also show cgo calls per second ?
		deltacalls := delta[uint64]()

		return func(_ time.Time, samples []metrics.Sample) any {
			calls := samples[idx_cgo_go_to_c_calls_calls].Value.Uint64()

			return []uint64{
				deltacalls(calls),
			}
		}
	},
	layout: Scatter{
		Name:   "cgo",
		Tags:   []tag{tagMisc},
		Title:  "CGO Calls",
		Type:   "bar",
		Events: "lastgc",
		Layout: ScatterLayout{
			BarMode: "stack",
			Yaxis: ScatterYAxis{
				Title: "calls",
			},
		},
		Subplots: []Subplot{
			{Name: "calls from go to c", Unitfmt: "%{y}"},
		},
		InfoText: `<i>calls from go to c</i> is the delta of <b>/cgo/go-to-c-calls:calls</b>. Count of calls made from Go to C by the current process.`,
	},
})
